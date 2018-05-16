package main

type Sender interface {
	Send(msg *Message, recipients []string) error
}
