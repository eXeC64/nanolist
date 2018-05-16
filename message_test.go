package main

import (
	"strings"
	"testing"
)

func TestParseMessage(t *testing.T) {
	fakeInput := strings.NewReader("" +
		"From: Example Person <test@example.com>\r\n" +
		"To: Alice <alice@example.com>, Bob <bob@example.com>\r\n" +
		"Cc: Charlie <charlie@example.com>, Dolores <dolores@example.com>\r\n" +
		"Bcc: Evan <evan@example.com>, Francis <francis@example.com>\r\n" +
		"Subject: My Test Subject\r\n" +
		"Date: Sat, 24 Nov 2035 11:45:15 −0500\r\n" +
		"Message-ID: test-msg-id@example.com\r\n" +
		"In-Reply-To: other-msg-id@example.com\r\n" +
		"\r\n" +
		"This is my body\nIt even has multiple lines.")

	msg, err := ParseMessage(fakeInput)
	if err != nil {
		t.Errorf("Parsing failed with error: %s", err.Error())
		return
	}

	if msg == nil {
		t.Errorf("Returned message was nil")
		return
	}

	if msg.From != "Example Person <test@example.com>" {
		t.Errorf("From parsed incorrectly")
	}

	if msg.To != "Alice <alice@example.com>, Bob <bob@example.com>" {
		t.Errorf("To parsed incorrectly")
	}

	if msg.Cc != "Charlie <charlie@example.com>, Dolores <dolores@example.com>" {
		t.Errorf("Cc parsed incorrectly")
	}

	if msg.Bcc != "Evan <evan@example.com>, Francis <francis@example.com>" {
		t.Errorf("Bcc parsed incorrectly")
	}

	if msg.Subject != "My Test Subject" {
		t.Errorf("Subject parsed incorrectly")
	}

	if msg.Date != "Sat, 24 Nov 2035 11:45:15 −0500" {
		t.Errorf("Date parsed incorrectly")
	}

	if msg.Id != "test-msg-id@example.com" {
		t.Errorf("Message id parsed incorrectly")
	}

	if msg.InReplyTo != "other-msg-id@example.com" {
		t.Errorf("In-Reply-To parsed incorrectly")
	}

	if msg.Body != "This is my body\nIt even has multiple lines." {
		t.Errorf("Body parsed incorrectly", msg.Body)
	}
}
