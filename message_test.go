package main

import (
	"net/mail"
	"strings"
	"testing"
	"time"
)

func TestParseMessage(t *testing.T) {
	fakeInput := strings.NewReader("" +
		"From: \"Example Person\" <test@example.com>\r\n" +
		"To: \"Alice\" <alice@example.com>, \"Bob\" <bob@example.com>\r\n" +
		"Cc: \"Charlie\" <charlie@example.com>, \"Dolores\" <dolores@example.com>\r\n" +
		"Bcc: \"Evan\" <evan@example.com>, \"Francis\" <francis@example.com>\r\n" +
		"Subject: My Test Subject\r\n" +
		"Date: Mon, 14 May 2018 19:41:34 +0200\r\n" +
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

	if msg.From.Name != "Example Person" {
		t.Errorf("From name parsed incorrectly")
	}

	if msg.From.Address != "test@example.com" {
		t.Errorf("From address parsed incorrectly")
	}

	if len(msg.To) != 2 {
		t.Errorf("To parsed with incorrect number of entries")
	} else {
		if msg.To[0].Name != "Alice" {
			t.Errorf("To first name parsed incorrectly")
		}
		if msg.To[0].Address != "alice@example.com" {
			t.Errorf("To first address parsed incorrectly")
		}

		if msg.To[1].Name != "Bob" {
			t.Errorf("To second name parsed incorrectly")
		}
		if msg.To[1].Address != "bob@example.com" {
			t.Errorf("To second address parsed incorrectly")
		}
	}

	if len(msg.Cc) != 2 {
		t.Errorf("Cc parsed with incorrect number of entries")
	} else {
		if msg.Cc[0].Name != "Charlie" {
			t.Errorf("Cc first name parsed incorrectly")
		}
		if msg.Cc[0].Address != "charlie@example.com" {
			t.Errorf("Cc first address parsed incorrectly")
		}

		if msg.Cc[1].Name != "Dolores" {
			t.Errorf("Cc second name parsed incorrectly")
		}
		if msg.Cc[1].Address != "dolores@example.com" {
			t.Errorf("Cc second address parsed incorrectly")
		}
	}

	if len(msg.Bcc) != 2 {
		t.Errorf("Bcc parsed with incorrect number of entries")
	} else {
		if msg.Bcc[0].Name != "Evan" {
			t.Errorf("Bcc first name parsed incorrectly")
		}
		if msg.Bcc[0].Address != "evan@example.com" {
			t.Errorf("Bcc first address parsed incorrectly")
		}

		if msg.Bcc[1].Name != "Francis" {
			t.Errorf("Bcc second name parsed incorrectly")
		}
		if msg.Bcc[1].Address != "francis@example.com" {
			t.Errorf("Bcc second address parsed incorrectly")
		}
	}

	if msg.Subject != "My Test Subject" {
		t.Errorf("Subject parsed incorrectly")
	}

	if msg.Date.String() != "2018-05-14 19:41:34 +0200 +0200" {
		t.Errorf("Date parsed incorrectly: %s", msg.Date.String())
	}

	if msg.Id != "test-msg-id@example.com" {
		t.Errorf("Message id parsed incorrectly")
	}

	if msg.InReplyTo != "other-msg-id@example.com" {
		t.Errorf("In-Reply-To parsed incorrectly")
	}

	if msg.Body != "This is my body\nIt even has multiple lines." {
		t.Error("Body parsed incorrectly")
	}
}

func TestParseSlimMessage(t *testing.T) {
	fakeInput := strings.NewReader("" +
		"From: \"Example Person\" <test@example.com>\r\n" +
		"To: \"Alice\" <alice@example.com>\r\n" +
		"Subject: My Test Subject\r\n" +
		"\r\n")

	msg, err := ParseMessage(fakeInput)
	if err != nil {
		t.Errorf("Parsing failed with error: %s", err.Error())
		return
	}

	if msg == nil {
		t.Errorf("Returned message was nil")
		return
	}

	if msg.From.Name != "Example Person" {
		t.Error("From name parsed incorrectly")
	}

	if msg.From.Address != "test@example.com" {
		t.Error("From name parsed incorrectly")
	}

	if msg.Subject != "My Test Subject" {
		t.Error("Subject parsed incorrectly")
	}

	if len(msg.To) != 1 {
		t.Error("To parsed incorrect number of addresses")
	} else {
		if msg.To[0].Name != "Alice" {
			t.Error("To name parsed incorrectly")
		}
		if msg.To[0].Address != "alice@example.com" {
			t.Error("To address parsed incorrectly")
		}
	}
}

func TestStringMessage(t *testing.T) {
	msg := &Message{
		Subject:   "Just a test subject",
		From:      &mail.Address{Name: "James Bond", Address: "bond@example.com"},
		Id:        "test-id@example.com",
		InReplyTo: "other-test-id@example.com",
		Body:      "This is my test body\nIt contains multiple lines!",
		To:        []*mail.Address{{"Alice", "alice@example.com"}, {"Bob", "bob@example.com"}},
		Cc:        []*mail.Address{{"Charlie", "charlie@example.com"}, {"Dolores", "dolores@example.com"}},
		Bcc:       []*mail.Address{{"Evan", "evan@example.com"}, {"Francis", "francis@example.com"}},
		Date:      time.Date(2018, time.May, 14, 19, 41, 34, 0, time.FixedZone("ABC", 2*60*60)),
	}

	str := msg.String()

	expected := "From: \"James Bond\" <bond@example.com>\r\n" +
		"To: \"Alice\" <alice@example.com>, \"Bob\" <bob@example.com>\r\n" +
		"Cc: \"Charlie\" <charlie@example.com>, \"Dolores\" <dolores@example.com>\r\n" +
		"Bcc: \"Evan\" <evan@example.com>, \"Francis\" <francis@example.com>\r\n" +
		"Date: Mon, 14 May 2018 19:41:34 +0200\r\n" +
		"Message-ID: test-id@example.com\r\n" +
		"In-Reply-To: other-test-id@example.com\r\n" +
		"Subject: Just a test subject\r\n" +
		"\r\n" +
		"This is my test body\nIt contains multiple lines!"

	if str != expected {
		t.Errorf("Message String-ified incorrectly.\nExpect: '%q'\n\n\nActual: '%q'", expected, str)
	}
}

func TestReply(t *testing.T) {
	msg := &Message{
		Subject: "This is a test",
		From:    &mail.Address{Name: "James Bond", Address: "bond@example.com"},
		Id:      "test-id@example.com",
	}

	reply := msg.Reply()

	if reply.Subject != "Re: This is a test" {
		t.Errorf("Incorrect Subject: %q", reply.Subject)
	}

	if len(reply.To) != 1 {
		t.Errorf("Incorrect number of recipients: %d", len(reply.To))
	} else {

		if reply.To[0].Name != "James Bond" {
			t.Errorf("Incorrect To name: %q", reply.To[0].Name)
		}

		if reply.To[0].Address != "bond@example.com" {
			t.Errorf("Incorrect To address: %q", reply.To[0].Address)
		}
	}

	if reply.InReplyTo != "test-id@example.com" {
		t.Errorf("Incorrect InReplyTo: %q", reply.InReplyTo)
	}
}

func TestAllRecipients(t *testing.T) {
	msg := &Message{
		To:  []*mail.Address{{"Charlie", "charlie@example.com"}, {"Dolores", "dolores@example.com"}},
		Cc:  []*mail.Address{{"Charles", "charlie@example.com"}, {"Bob", "bob@example.com"}},
		Bcc: []*mail.Address{{"Evan", "evan@example.com"}, {"Francis", "francis@example.com"}},
	}

	addrs := AddressesOnly(msg.AllRecipients())

	expected := []string{"bob@example.com",
		"charlie@example.com",
		"dolores@example.com",
		"evan@example.com",
		"francis@example.com",
	}

	if len(expected) != len(addrs) {
		t.Errorf("Incorrect number of addresses: %d", len(addrs))
		return
	}

	for i := range addrs {
		if addrs[i] != expected[i] {
			t.Errorf("Incorrect address: %q != %q", addrs[i], expected[i])
		}
	}
}
