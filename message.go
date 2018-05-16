package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"time"
)

type Message struct {
	Subject     string
	From        *mail.Address
	To          []*mail.Address
	Cc          []*mail.Address
	Bcc         []*mail.Address
	Date        time.Time
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

	date, err := inMessage.Header.Date()
	if err != nil {
		return nil, err
	}

	from, err := mail.ParseAddress(inMessage.Header.Get("From"))
	if err != nil {
		return nil, err
	}

	to, err := inMessage.Header.AddressList("To")
	if err != nil {
		return nil, err
	}

	cc, err := inMessage.Header.AddressList("Cc")
	if err != nil {
		return nil, err
	}

	bcc, err := inMessage.Header.AddressList("Bcc")
	if err != nil {
		return nil, err
	}

	msg := &Message{
		Subject:   inMessage.Header.Get("Subject"),
		From:      from,
		Id:        inMessage.Header.Get("Message-ID"),
		InReplyTo: inMessage.Header.Get("In-Reply-To"),
		Body:      string(body[:]),
		To:        to,
		Cc:        cc,
		Bcc:       bcc,
		Date:      date,
	}

	return msg, nil
}

func (msg *Message) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "From: %s\r\n", msg.From)
	fmt.Fprintf(&buf, "To: %s\r\n", msg.To)
	fmt.Fprintf(&buf, "Cc: %s\r\n", msg.Cc)
	fmt.Fprintf(&buf, "Bcc: %s\r\n", msg.Bcc)
	if !msg.Date.IsZero() {
		fmt.Fprintf(&buf, "Date: %s\r\n", msg.Date.Format("Mon, 2 Jan 2006 15:04:05 -0700"))
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
