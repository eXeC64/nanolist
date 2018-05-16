package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
)

type Message struct {
	Subject     string
	From        string
	To          string
	Cc          string
	Bcc         string
	Date        string
	Id          string
	InReplyTo   string
	ContentType string
	XList       string
	Body        string
}

func ParseMessage(input io.Reader) (*Message, error) {

	inMessage, err := mail.ReadMessage(input)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(inMessage.Body)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Subject:   inMessage.Header.Get("Subject"),
		From:      inMessage.Header.Get("From"),
		Id:        inMessage.Header.Get("Message-ID"),
		InReplyTo: inMessage.Header.Get("In-Reply-To"),
		Body:      string(body[:]),
		To:        inMessage.Header.Get("To"),
		Cc:        inMessage.Header.Get("Cc"),
		Bcc:       inMessage.Header.Get("Bcc"),
		Date:      inMessage.Header.Get("Date"),
	}

	return msg, nil
}

func (msg *Message) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "From: %s\r\n", msg.From)
	fmt.Fprintf(&buf, "To: %s\r\n", msg.To)
	fmt.Fprintf(&buf, "Cc: %s\r\n", msg.Cc)
	fmt.Fprintf(&buf, "Bcc: %s\r\n", msg.Bcc)
	if len(msg.Date) > 0 {
		fmt.Fprintf(&buf, "Date: %s\r\n", msg.Date)
	}
	if len(msg.Id) > 0 {
		fmt.Fprintf(&buf, "Messsage-ID: %s\r\n", msg.Id)
	}
	fmt.Fprintf(&buf, "In-Reply-To: %s\r\n", msg.InReplyTo)
	if len(msg.XList) > 0 {
		fmt.Fprintf(&buf, "X-Mailing-List: %s\r\n", msg.XList)
		fmt.Fprintf(&buf, "List-ID: %s\r\n", msg.XList)
		fmt.Fprintf(&buf, "Sender: %s\r\n", msg.XList)
	}
	if len(msg.ContentType) > 0 {
		fmt.Fprintf(&buf, "Content-Type: %s\r\n", msg.ContentType)
	}
	fmt.Fprintf(&buf, "Subject: %s\r\n", msg.Subject)
	fmt.Fprintf(&buf, "\r\n%s", msg.Body)

	return buf.String()
}
