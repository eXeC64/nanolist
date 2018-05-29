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
const noSuchList string = "No such list exists. Please check you entered its address correctly."

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

	// Find out if it's for any of our mailing lists
	recipients := AddressesOnly(msg.AllRecipients())
	allLists, err := p.Lists.FetchListAddresses()
	if err != nil {
		p.Log.Printf("Failed to fetch list addresses: %q", err.Error())
		p.sendReply(msg, errMsg)
		return
	}

	// Both lists are sorted, so get the intersection of them to find the lists
	// we need to send to.
	toLists := IntersectionOf(recipients, allLists)

	if len(toLists) == 0 {
		p.sendReply(msg, "No mailing lists addressed. Your message has not been delivered.")
		return
	}

	for _, listAddr := range toLists {
		list, err := p.Lists.FetchList(listAddr)
		if err != nil {
			p.Log.Printf("Failed to fetch list: %q", err.Error())
			p.sendReply(msg, errMsg)
			return
		}

		// If there's a whitelist of posters, check the sender is on it
		if len(list.Posters) > 0 {
			isPoster := false
			for _, poster := range list.Posters {
				if msg.From.Address == poster {
					isPoster = true
					break
				}
			}
			if !isPoster {
				p.sendReply(msg, "You are not permitted to post to "+list.Address)
				return
			}
		} else if list.SubscribersOnly {
			// If this is a subscribers-only list, check the sender is subscribed
			isSubscribed, err := p.Subscriptions.IsSubscribed(msg.From.Address, list.Address)
			if err != nil {
				p.Log.Printf("Failed to determine if user is subscribed: %q", err.Error())
				p.sendReply(msg, errMsg)
				return
			}

			if !isSubscribed {
				p.sendReply(msg, "Only subscribers may post to "+list.Address)
				return
			}
		}

		p.sendToList(msg, list)
	}
}

func (p *Postman) sendToList(msg *Message, list *List) {
	// Construct a version of the message to relay to subscribers of this list
	listMsg := msg.RelayVia(list.Address)

	// Fetch the subscribers of this list
	recipients, err := p.Subscriptions.FetchSubscribers(list.Address)
	if err != nil {
		p.Log.Printf("Failed to fetch subscribers: %q", err.Error())
		return
	}

	// Add list specific bcc's to this message
	recipients = append(recipients, list.Bcc...)

	// Relay the message
	p.Sender.Send(listMsg, recipients)
}

func (p *Postman) sendReply(msg *Message, response string) {
	reply := msg.Reply()
	reply.From = &mail.Address{Name: "", Address: p.CommandAddress}
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
    subscribe <list-address>: Subscribe to the given mailing list
    unsubscribe <list-address>: Unsubscribe from the given mailing list`)
}

func (p *Postman) handleListsCommand(msg *Message) {
	addrs, err := p.Lists.FetchListAddresses()
	if err != nil {
		p.Log.Printf("Failed to fetch list addresses: %q", err.Error())
		p.sendReply(msg, errMsg)
		return
	}

	sort.Strings(addrs)
	lists := []*List{}
	for _, address := range addrs {
		list, err := p.Lists.FetchList(address)
		if err != nil {
			p.Log.Printf("Failed to fetch list address %q: %q", address, err.Error())
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
				"Name: %s\r\n"+
					"Description: %s\r\n"+
					"Address: %s\r\n\r\n",
				list.Name, list.Description, list.Address)
		}
	}

	fmt.Fprintf(&body,
		"\r\nTo subscribe to a mailing list, email %s with 'subscribe <list-address>' as the subject.\r\n",
		p.CommandAddress)

	p.sendReply(msg, body.String())
}

func (p *Postman) handleSubscribeCommand(msg *Message, args []string) {
	if len(args) < 1 {
		p.sendReply(msg, "No mailing list address specified. Unable to subscribe you.")
		return
	}

	listAddr := args[0]

	exists, err := p.Lists.IsValidList(listAddr)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}
	if !exists {
		p.sendReply(msg, noSuchList)
		return
	}

	isSubscribed, err := p.Subscriptions.IsSubscribed(msg.From.Address, listAddr)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}
	if isSubscribed {
		p.sendReply(msg, "You are already subscribed to "+listAddr)
		return
	}

	err = p.Subscriptions.Subscribe(msg.From.Address, listAddr)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}

	p.sendReply(msg, "You have been subscribed to "+listAddr)
}

func (p *Postman) handleUnsubscribeCommand(msg *Message, args []string) {
	if len(args) < 1 {
		p.sendReply(msg, "No mailing list address specified. Unable to unsubscribe you.")
		return
	}

	listAddr := args[0]

	exists, err := p.Lists.IsValidList(listAddr)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}
	if !exists {
		p.sendReply(msg, noSuchList)
		return
	}

	isSubscribed, err := p.Subscriptions.IsSubscribed(msg.From.Address, listAddr)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}
	if !isSubscribed {
		p.sendReply(msg, "You are not subscribed to "+listAddr)
		return
	}

	err = p.Subscriptions.Unsubscribe(msg.From.Address, listAddr)
	if err != nil {
		p.sendReply(msg, errMsg)
		return
	}

	p.sendReply(msg, "You have been unsubscribed from "+listAddr)
}
