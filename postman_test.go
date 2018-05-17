package main

import (
	"github.com/stretchr/testify/mock"
	"log"
	"strings"
	"testing"
)

// Writer that discards all input
type NullWriter struct {
}

func (w *NullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func TestHelpCommand(t *testing.T) {

	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          &MemoryListManager{},
	}

	senderMock.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	input := strings.NewReader("To: test@example.com\r\n" +
		"From: user@example.com\r\n" +
		"Subject: help\r\n" +
		"\r\n\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	if len(senderMock.Calls) < 1 {
		t.Errorf("Send not called")
		return
	}

	args := senderMock.Calls[0].Arguments

	recipients := args.Get(1).([]string)
	if len(recipients) != 1 {
		t.Errorf("Wrong number of recipients: %d", len(recipients))
	} else {
		addr := recipients[0]
		if addr != "user@example.com" {
			t.Errorf("Wrong recipient: %s", addr)
		}
	}

	msg := args.Get(0).(*Message)

	if !strings.Contains(msg.Body, "help: Reply with this help information") {
		t.Errorf("Response body did not contain expected help information: %q", msg.Body)
	}

	if msg.To[0].Address != "user@example.com" {
		t.Errorf("Response To address was incorrect: %q", msg.To[0].Address)
	}
}
