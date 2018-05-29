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

func checkResponse(t *testing.T, sender *MockSender, to string, body string) {

	if len(sender.Calls) < 1 {
		t.Errorf("Send not called")
	} else {

		args := sender.Calls[0].Arguments

		recipients := args.Get(1).([]string)
		if len(recipients) != 1 {
			t.Errorf("Wrong number of recipients: %d", len(recipients))
		} else {
			addr := recipients[0]
			if addr != to {
				t.Errorf("Wrong recipient: %q Expcted: %q", addr, to)
			}
		}

		msg := args.Get(0).(*Message)

		if !strings.Contains(msg.Body, body) {
			t.Errorf("Response body did not match expectations. Body: %q Expected: %q", msg.Body, body)
		}

		if msg.To[0].Address != to {
			t.Errorf("Response To address was incorrect: To: %q Expected: %q", msg.To[0].Address, to)
		}
	}
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

	checkResponse(t, senderMock, "user@example.com", "help: Reply with this help information")
}

func TestListsCommand(t *testing.T) {

	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}
	listManager := &MemoryListManager{}

	listManager.Add(&List{
		Name:        "Poker Discussion",
		Description: "All things poker",
		Address:     "poker@example.com",
	})

	listManager.Add(&List{
		Name:        "Secret Chat",
		Description: "Sssh",
		Address:     "secret@example.com",
		Hidden:      true,
	})

	listManager.Add(&List{
		Name:        "Nomic",
		Description: "Lets play nomic",
		Address:     "nomic-business@example.com",
	})

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          listManager,
	}

	senderMock.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	input := strings.NewReader("To: test@example.com\r\n" +
		"From: user@example.com\r\n" +
		"Subject: lists\r\n" +
		"\r\n\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	checkResponse(t, senderMock, "user@example.com", "Name: Nomic\r\n"+
		"Description: Lets play nomic\r\n"+
		"Address: nomic-business@example.com\r\n"+
		"\r\n"+
		"Name: Poker Discussion\r\n"+
		"Description: All things poker\r\n"+
		"Address: poker@example.com\r\n")
}

func TestSubscribeCommand(t *testing.T) {

	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}
	listManager := &MemoryListManager{}

	listManager.Add(&List{
		Name:        "Poker Discussion",
		Description: "All things poker",
		Address:     "poker@example.com",
	})

	listManager.Add(&List{
		Name:        "Nomic",
		Description: "Lets play nomic",
		Address:     "nomic-business@example.com",
	})

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          listManager,
	}

	subManagerMock.On("IsSubscribed", mock.Anything, mock.Anything).Return(false, nil).Once()
	subManagerMock.On("Subscribe", mock.Anything, mock.Anything).Return(nil).Once()
	senderMock.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	input := strings.NewReader("To: test@example.com\r\n" +
		"From: user@example.com\r\n" +
		"Subject: subscribe nomic-business@example.com\r\n" +
		"\r\n\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	if len(subManagerMock.Calls) < 2 {
		t.Errorf("Subscribe or IsSubscribed not called")
	} else {

		// IsSubscribed call
		args := subManagerMock.Calls[0].Arguments
		if args.String(0) != "user@example.com" {
			t.Errorf("Incorrect email used in IsSubscribed call")
		}
		if args.String(1) != "nomic-business@example.com" {
			t.Errorf("Incorrect list address used in IsSubscribed call")
		}

		// Subscribe call
		args = subManagerMock.Calls[1].Arguments
		if args.String(0) != "user@example.com" {
			t.Errorf("Incorrect email used in Subscribe call")
		}
		if args.String(1) != "nomic-business@example.com" {
			t.Errorf("Incorrect list address used in Subscribe call")
		}
	}

	// Send call
	checkResponse(t, senderMock, "user@example.com", "You have been subscribed to nomic-business@example.com")
}

func TestUnsubscribeCommand(t *testing.T) {

	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}
	listManager := &MemoryListManager{}

	listManager.Add(&List{
		Name:        "Poker Discussion",
		Description: "All things poker",
		Address:     "poker@example.com",
	})

	listManager.Add(&List{
		Name:        "Nomic",
		Description: "Lets play nomic",
		Address:     "nomic-business@example.com",
	})

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          listManager,
	}

	subManagerMock.On("IsSubscribed", mock.Anything, mock.Anything).Return(true, nil).Once()
	subManagerMock.On("Unsubscribe", mock.Anything, mock.Anything).Return(nil).Once()
	senderMock.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	input := strings.NewReader("To: test@example.com\r\n" +
		"From: user@example.com\r\n" +
		"Subject: unsubscribe nomic-business@example.com\r\n" +
		"\r\n\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	if len(subManagerMock.Calls) < 2 {
		t.Errorf("Unsubscribe or IsSubscribed not called")
	} else {

		// IsSubscribed call
		args := subManagerMock.Calls[0].Arguments
		if args.String(0) != "user@example.com" {
			t.Errorf("Incorrect email used in IsSubscribed call")
		}
		if args.String(1) != "nomic-business@example.com" {
			t.Errorf("Incorrect list address used in IsSubscribed call")
		}

		// Unsubscribe call
		args = subManagerMock.Calls[1].Arguments
		if args.String(0) != "user@example.com" {
			t.Errorf("Incorrect email used in Unsubscribe call")
		}
		if args.String(1) != "nomic-business@example.com" {
			t.Errorf("Incorrect list address used in Unsubscribe call")
		}
	}

	// Send call
	checkResponse(t, senderMock, "user@example.com", "You have been unsubscribed from nomic-business@example.com")
}

func TestRelayMessage(t *testing.T) {
	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}
	listManager := &MemoryListManager{}

	listManager.Add(&List{
		Name:        "Poker Discussion",
		Description: "All things poker",
		Address:     "poker@example.com",
	})

	listManager.Add(&List{
		Name:        "Nomic",
		Description: "Lets play nomic",
		Address:     "nomic-business@example.com",
	})

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          listManager,
	}

	subs := []string{
		"user@example.com",
		"admin@example.com",
	}

	subManagerMock.On("FetchSubscribers", "poker@example.com").Return(subs, nil).Once()
	senderMock.On("Send", mock.Anything, subs).Return(nil).Once()

	input := strings.NewReader("To: poker@example.com\r\n" +
		"From: user@example.com\r\n" +
		"Subject: example message\r\n" +
		"\r\n" +
		"Hello, this is my message." +
		"\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	if len(subManagerMock.Calls) < 1 {
		t.Errorf("FetchSubscribers not called")
	}

	if len(senderMock.Calls) < 1 {
		t.Errorf("Send not called")
	} else {
		// Check Send call
		msg := senderMock.Calls[0].Arguments.Get(0).(*Message)
		actual := msg.String()
		expected := "From: <user@example.com>\r\n" +
			"To: <poker@example.com>\r\n" +
			"X-Mailing-List: poker@example.com\r\n" +
			"List-ID: poker@example.com\r\n" +
			"Sender: poker@example.com\r\n" +
			"Subject: example message\r\n" +
			"\r\n" +
			"Hello, this is my message." +
			"\r\n"

		if expected != actual {
			t.Errorf("Expected message != actual message.\nExpected:\n%q\nActual:\n%q", expected, actual)
		}
	}
}

func TestDontRelaySubscribersOnlyMessage(t *testing.T) {
	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}
	listManager := &MemoryListManager{}

	listManager.Add(&List{
		Name:            "Poker Discussion",
		Description:     "All things poker",
		Address:         "poker@example.com",
		SubscribersOnly: true,
	})

	listManager.Add(&List{
		Name:        "Nomic",
		Description: "Lets play nomic",
		Address:     "nomic-business@example.com",
	})

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          listManager,
	}

	subs := []string{
		"user@example.com",
		"admin@example.com",
	}

	subManagerMock.On("FetchSubscribers", "poker@example.com").Return(subs, nil).Once()
	subManagerMock.On("IsSubscribed", "spammer@example.com", "poker@example.com").Return(false, nil).Once()
	senderMock.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	input := strings.NewReader("To: poker@example.com\r\n" +
		"From: spammer@example.com\r\n" +
		"Subject: example message\r\n" +
		"\r\n" +
		"Hello, this is my message." +
		"\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	if len(subManagerMock.Calls) < 1 {
		t.Errorf("FetchSubscribers not called")
	}

	checkResponse(t, senderMock, "spammer@example.com", "Only subscribers may post to poker@example.com")
}

func TestDontRelayPostersOnlyMessage(t *testing.T) {
	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}
	listManager := &MemoryListManager{}

	listManager.Add(&List{
		Name:        "Poker Discussion",
		Description: "All things poker",
		Address:     "poker@example.com",
		Posters:     []string{"admin@example.com"},
	})

	listManager.Add(&List{
		Name:        "Nomic",
		Description: "Lets play nomic",
		Address:     "nomic-business@example.com",
	})

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          listManager,
	}

	senderMock.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	input := strings.NewReader("To: poker@example.com\r\n" +
		"From: spammer@example.com\r\n" +
		"Subject: example message\r\n" +
		"\r\n" +
		"Hello, this is my message." +
		"\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	checkResponse(t, senderMock, "spammer@example.com", "You are not permitted to post to poker@example.com")
}

func TestRelayMessageToListBccs(t *testing.T) {
	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}
	listManager := &MemoryListManager{}

	listManager.Add(&List{
		Name:        "Poker Discussion",
		Description: "All things poker",
		Address:     "poker@example.com",
		Bcc:         []string{"cheater@example.com", "shark@example.com"},
	})

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          listManager,
	}

	subs := []string{
		"user@example.com",
		"admin@example.com",
	}

	recipients := []string{
		"user@example.com",
		"admin@example.com",
		"cheater@example.com",
		"shark@example.com",
	}

	subManagerMock.On("FetchSubscribers", "poker@example.com").Return(subs, nil).Once()
	senderMock.On("Send", mock.Anything, recipients).Return(nil).Once()

	input := strings.NewReader("To: poker@example.com\r\n" +
		"From: user@example.com\r\n" +
		"Subject: example message\r\n" +
		"\r\n" +
		"Hello, this is my message." +
		"\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	if len(subManagerMock.Calls) < 1 {
		t.Errorf("FetchSubscribers not called")
	}

	if len(senderMock.Calls) < 1 {
		t.Errorf("Send not called")
	}
}

func TestRelayNowhereErrors(t *testing.T) {
	// GIVEN
	senderMock := &MockSender{}
	subManagerMock := &MockSubscriptionManager{}
	listManager := &MemoryListManager{}

	pm := &Postman{
		CommandAddress: "test@example.com",
		Log:            log.New(&NullWriter{}, "", 0),
		Sender:         senderMock,
		Subscriptions:  subManagerMock,
		Lists:          listManager,
	}

	senderMock.On("Send", mock.Anything, mock.Anything).Return(nil).Once()

	input := strings.NewReader("To: poker@example.com\r\n" +
		"From: user@example.com\r\n" +
		"Subject: example message\r\n" +
		"\r\n" +
		"Hello, this is my message." +
		"\r\n")

	// WHEN
	pm.HandleMail(input)

	// THEN
	checkResponse(t, senderMock, "user@example.com", "No mailing lists addressed. Your message has not been delivered.")
}
