package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"sort"
	"strings"
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

	// Pull these fields if available
	date, _ := inMessage.Header.Date()
	from, _ := mail.ParseAddress(inMessage.Header.Get("From"))
	to, _ := inMessage.Header.AddressList("To")
	cc, _ := inMessage.Header.AddressList("Cc")
	bcc, _ := inMessage.Header.AddressList("Bcc")

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

func formatAddressList(addresses []*mail.Address) string {
	strs := []string{}

	for _, addr := range addresses {
		strs = append(strs, addr.String())
	}

	return strings.Join(strs, ", ")
}

func (msg *Message) String() string {
	var buf bytes.Buffer

	fmt.Fprintf(&buf, "From: %s\r\n", msg.From)
	fmt.Fprintf(&buf, "To: %s\r\n", formatAddressList(msg.To))
	fmt.Fprintf(&buf, "Cc: %s\r\n", formatAddressList(msg.Cc))
	fmt.Fprintf(&buf, "Bcc: %s\r\n", formatAddressList(msg.Bcc))
	if !msg.Date.IsZero() {
		fmt.Fprintf(&buf, "Date: %s\r\n", msg.Date.Format("Mon, 2 Jan 2006 15:04:05 -0700"))
	}
	if len(msg.Id) > 0 {
		fmt.Fprintf(&buf, "Message-ID: %s\r\n", msg.Id)
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

func (msg *Message) Reply() *Message {
	return &Message{
		Subject:   "Re: " + msg.Subject,
		To:        []*mail.Address{msg.From},
		Date:      time.Now(),
		InReplyTo: msg.Id,
	}
}

func (msg *Message) AllRecipients() []*mail.Address {
	addrs := []*mail.Address{}

	for _, to := range msg.To {
		addrs = append(addrs, to)
	}
	for _, cc := range msg.Cc {
		addrs = append(addrs, cc)
	}
	for _, bcc := range msg.Bcc {
		addrs = append(addrs, bcc)
	}

	sort.Slice(addrs, func(i, j int) bool {
		return addrs[i].Address < addrs[j].Address
	})

	// Remove duplicates, going back to front
	for i := len(addrs) - 1; i > 0; i-- {
		if addrs[i].Address == addrs[i-1].Address {
			addrs = append(addrs[:i], addrs[i+1:]...)
		}
	}

	return addrs
}
