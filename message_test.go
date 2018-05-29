package main

import (
	"net/mail"
	"strings"
	"testing"
	"time"
)

func checkMessage(expected *Message, actual *Message, t *testing.T) {

	if expected.From != nil {
		if expected.From.Name != actual.From.Name {
			t.Errorf("Expected from name (%q) != actual from name %q", expected.From.Name, actual.From.Name)
		}

		if expected.From.Address != actual.From.Address {
			t.Errorf("Expected address (%q) != actual address %q", expected.From.Address, actual.From.Address)
		}
	}

	if len(expected.To) != len(actual.To) {
		t.Errorf("Expected To length (%d) != Actual To length (%d)", len(expected.To), len(actual.To))
	} else {
		for i := range expected.To {
			if expected.To[i].Name != actual.To[i].Name {
				t.Errorf("Expected To[%d] name (%q) != Actual To[%d] name (%q)", i, expected.To[i].Name, i, actual.To[i].Name)
			}
			if expected.To[i].Address != actual.To[i].Address {
				t.Errorf("Expected To[%d] address (%q) != Actual To[%d] address (%q)", i, expected.To[i].Address, i, actual.To[i].Address)
			}
		}
	}

	if len(expected.Cc) != len(actual.Cc) {
		t.Errorf("Expected Cc length (%d) != Actual Cc length (%d)", len(expected.Cc), len(actual.Cc))
	} else {
		for i := range expected.Cc {
			if expected.Cc[i].Name != actual.Cc[i].Name {
				t.Errorf("Expected Cc[%d] name (%q) != Actual Cc[%d] name (%q)", i, expected.Cc[i].Name, i, actual.Cc[i].Name)
			}
			if expected.Cc[i].Address != actual.Cc[i].Address {
				t.Errorf("Expected Cc[%d] address (%q) != Actual Cc[%d] address (%q)", i, expected.Cc[i].Address, i, actual.Cc[i].Address)
			}
		}
	}

	if len(expected.Bcc) != len(actual.Bcc) {
		t.Errorf("Expected Bcc length (%d) != Actual Bcc length (%d)", len(expected.Bcc), len(actual.Bcc))
	} else {
		for i := range expected.Bcc {
			if expected.Bcc[i].Name != actual.Bcc[i].Name {
				t.Errorf("Expected Bcc[%d] name (%q) != Actual Bcc[%d] name (%q)", i, expected.Bcc[i].Name, i, actual.Bcc[i].Name)
			}
			if expected.Bcc[i].Address != actual.Bcc[i].Address {
				t.Errorf("Expected Bcc[%d] address (%q) != Actual Bcc[%d] address (%q)", i, expected.Bcc[i].Address, i, actual.Bcc[i].Address)
			}
		}
	}

	if expected.Subject != actual.Subject {
		t.Errorf("Expected Subject (%q) != Actual Subject (%q)", expected.Subject, actual.Subject)
	}

	if expected.Date.String() != actual.Date.String() {
		t.Errorf("Expected Date (%q) != Actual Date (%q)", expected.Date.String(), actual.Date.String())
	}

	if expected.Id != actual.Id {
		t.Errorf("Expected Id (%q) != Actual Id (%q)", expected.Id, actual.Id)
	}

	if expected.InReplyTo != actual.InReplyTo {
		t.Errorf("Expected InReplyTo (%q) != Actual InReplyTo (%q)", expected.InReplyTo, actual.InReplyTo)
	}

	if expected.ContentType != actual.ContentType {
		t.Errorf("Expected ContentType (%q) != Actual ContentType (%q)", expected.ContentType, actual.ContentType)
	}

	if expected.Body != actual.Body {
		t.Errorf("Expected Body (%q) != Actual Body (%q)", expected.Body, actual.Body)
	}
}

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
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"This is my body\nIt even has multiple lines.")

	actual, err := ParseMessage(fakeInput)
	if err != nil {
		t.Errorf("Parsing failed with error: %s", err.Error())
		return
	}

	if actual == nil {
		t.Errorf("Returned message was nil")
		return
	}

	expected := &Message{
		Subject:     "My Test Subject",
		From:        &mail.Address{Name: "Example Person", Address: "test@example.com"},
		Id:          "test-msg-id@example.com",
		InReplyTo:   "other-msg-id@example.com",
		ContentType: "text/plain",
		Body:        "This is my body\nIt even has multiple lines.",
		To:          []*mail.Address{{"Alice", "alice@example.com"}, {"Bob", "bob@example.com"}},
		Cc:          []*mail.Address{{"Charlie", "charlie@example.com"}, {"Dolores", "dolores@example.com"}},
		Bcc:         []*mail.Address{{"Evan", "evan@example.com"}, {"Francis", "francis@example.com"}},
		Date:        time.Date(2018, time.May, 14, 19, 41, 34, 0, time.FixedZone("+0200", 2*60*60)),
	}

	checkMessage(expected, actual, t)
}

func TestParseSlimMessage(t *testing.T) {
	fakeInput := strings.NewReader("" +
		"From: \"Example Person\" <test@example.com>\r\n" +
		"To: \"Alice\" <alice@example.com>\r\n" +
		"Subject: My Test Subject\r\n" +
		"\r\n")

	actual, err := ParseMessage(fakeInput)
	if err != nil {
		t.Errorf("Parsing failed with error: %s", err.Error())
		return
	}

	if actual == nil {
		t.Errorf("Returned message was nil")
		return
	}

	expected := &Message{
		Subject: "My Test Subject",
		From:    &mail.Address{Name: "Example Person", Address: "test@example.com"},
		To:      []*mail.Address{{"Alice", "alice@example.com"}},
	}

	checkMessage(expected, actual, t)
}

func TestStringMessage(t *testing.T) {
	msg := &Message{
		Subject:     "Just a test subject",
		From:        &mail.Address{Name: "James Bond", Address: "bond@example.com"},
		Id:          "test-id@example.com",
		InReplyTo:   "other-test-id@example.com",
		ContentType: "text/plain",
		Body:        "This is my test body\nIt contains multiple lines!",
		To:          []*mail.Address{{"Alice", "alice@example.com"}, {"Bob", "bob@example.com"}},
		Cc:          []*mail.Address{{"Charlie", "charlie@example.com"}, {"Dolores", "dolores@example.com"}},
		Bcc:         []*mail.Address{{"Evan", "evan@example.com"}, {"Francis", "francis@example.com"}},
		Date:        time.Date(2018, time.May, 14, 19, 41, 34, 0, time.FixedZone("ABC", 2*60*60)),
	}

	str := msg.String()

	expected := "From: \"James Bond\" <bond@example.com>\r\n" +
		"To: \"Alice\" <alice@example.com>, \"Bob\" <bob@example.com>\r\n" +
		"Cc: \"Charlie\" <charlie@example.com>, \"Dolores\" <dolores@example.com>\r\n" +
		"Bcc: \"Evan\" <evan@example.com>, \"Francis\" <francis@example.com>\r\n" +
		"Date: Mon, 14 May 2018 19:41:34 +0200\r\n" +
		"Message-ID: test-id@example.com\r\n" +
		"In-Reply-To: other-test-id@example.com\r\n" +
		"Content-Type: text/plain\r\n" +
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

	expected := &Message{
		Subject:   "Re: This is a test",
		To:        []*mail.Address{{"James Bond", "bond@example.com"}},
		InReplyTo: "test-id@example.com",
		Date:      time.Now(),
	}

	actual := msg.Reply()

	// Cheat this one field, because time passes
	actual.Date = expected.Date

	checkMessage(expected, actual, t)
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
