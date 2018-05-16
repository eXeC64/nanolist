package main

import (
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
