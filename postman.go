package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/mail"
	"sort"
	"strings"
)

type Postman struct {
	CommandAddress string
	Log            *log.Logger
	Sender         Sender
	Subscriptions  SubscriptionManager
	Lists          ListManager
}

const errMsg string = "There was an internal error. Please try again later."

func (p *Postman) HandleMail(input io.Reader) {

	msg, err := ParseMessage(input)
	if err != nil {
		p.Log.Printf("Failed to parse message: %q", err.Error())
		return
	}

	// Check if it's to the command address - but only direct To's count
	isCommand := false
	for _, addr := range msg.To {
		if addr.Address == p.CommandAddress {
			isCommand = true
			break
		}
	}

	// Intended for the command address, handle that
	if isCommand {
		p.handleCommand(msg)
		return
	}

	// Find out if it's for any mailing lists
	// TODO
}

func (p *Postman) sendReply(msg *Message, response string) {
	reply := msg.Reply()
	reply.From = &mail.Address{"", p.CommandAddress}
	reply.Body = response
	p.Sender.Send(reply, []string{reply.To[0].Address})
}

func (p *Postman) handleCommand(msg *Message) {
	parts := strings.Split(msg.Subject, " ")
	if len(parts) < 1 {
		p.sendReply(msg, "No command specified")
		return
	}

	cmd := strings.ToLower(parts[0])
	args := parts[1:]

	if cmd == "help" {
		p.handleHelpCommand(msg)
	} else if cmd == "lists" {
		p.handleListsCommand(msg)
	} else if cmd == "subscribe" {
		p.handleSubscribeCommand(msg, args)
	} else if cmd == "unsubscribe" {
		p.handleUnsubscribeCommand(msg, args)
	} else {
		p.sendReply(msg, "Unknown command")
	}
}

func (p *Postman) handleHelpCommand(msg *Message) {
	p.sendReply(msg, `Available commands:
    help: Reply with this help information
    lists: Reply with list of available mailing lists
    subscribe <list-id>: Subscribe to the given mailing list
    unsubscribe <list-id>: Unsubscribe from the given mailing list`)
}

func (p *Postman) handleListsCommand(msg *Message) {
	ids, err := p.Lists.FetchListIds()
	if err != nil {
		p.Log.Printf("Failed to fetch list ids: %q", err.Error())
		p.sendReply(msg, errMsg)
		return
	}

	sort.Strings(ids)
	lists := []*List{}
	for _, id := range ids {
		list, err := p.Lists.FetchList(id)
		if err != nil {
			p.Log.Printf("Failed to fetch list id %q: %q", id, err.Error())
			p.sendReply(msg, errMsg)
			return
		}
		lists = append(lists, list)
	}

	var body bytes.Buffer
	fmt.Fprintf(&body, "Available mailing lists:\r\n\r\n")
	for _, list := range lists {
		if !list.Hidden {
			fmt.Fprintf(&body,
				"Id: %s\r\n"+
					"Name: %s\r\n"+
					"Description: %s\r\n"+
					"Address: %s\r\n\r\n",
				list.Id, list.Name, list.Description, list.Address)
		}
	}

	fmt.Fprintf(&body,
		"\r\nTo subscribe to a mailing list, email %s with 'subscribe <list-id>' as the subject.\r\n",
		p.CommandAddress)

	p.sendReply(msg, body.String())
}

func (p *Postman) handleSubscribeCommand(msg *Message, args []string) {
	if len(args) < 1 {
		p.sendReply(msg, "No mailing list id specified. Unable to subscribe you.")
		return
	}

	const noSuchList string = "No such list exists. Please check you entered its id correctly."

	listId := args[0]

	exists, err := p.Lists.IsValidList(listId)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}
	if !exists {
		p.sendReply(msg, noSuchList)
		return
	}

	isSubscribed, err := p.Subscriptions.IsSubscribed(msg.From.Address, listId)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}
	if isSubscribed {
		p.sendReply(msg, "You are already subscribed to "+listId)
		return
	}

	err = p.Subscriptions.Subscribe(msg.From.Address, listId)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}

	p.sendReply(msg, "You have been subscribed to "+listId)
}

func (p *Postman) handleUnsubscribeCommand(msg *Message, args []string) {
	if len(args) < 1 {
		p.sendReply(msg, "No mailing list id specified. Unable to unsubscribe you.")
		return
	}

	const noSuchList string = "No such list exists. Please check you entered its id correctly."

	listId := args[0]

	exists, err := p.Lists.IsValidList(listId)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}
	if !exists {
		p.sendReply(msg, noSuchList)
		return
	}

	isSubscribed, err := p.Subscriptions.IsSubscribed(msg.From.Address, listId)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}
	if !isSubscribed {
		p.sendReply(msg, "You are not subscribed to "+listId)
		return
	}

	err = p.Subscriptions.Unsubscribe(msg.From.Address, listId)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}

	p.sendReply(msg, "You have been unsubscribed from "+listId)
}
