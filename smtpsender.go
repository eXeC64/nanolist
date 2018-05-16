package main

import (
	"errors"
	"net/smtp"
	"strconv"
)

type SMTPSender struct {
	host     string
	port     string
	username string
	password string
	ready    bool
}

func (s *SMTPSender) Login(host string, port int, username string, password string) error {
	s.host = host
	s.port = strconv.Itoa(port)
	s.username = username
	s.password = password
	// Todo - validate auth details
	s.ready = true
	return nil
}

func (s *SMTPSender) Send(msg *Message, recipients []string) error {
	if !s.ready {
		return errors.New("SMTPSender.Send called before ready")
	}
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	return smtp.SendMail(s.host+":"+s.port, auth, msg.From, recipients, []byte(msg.String()))
}
