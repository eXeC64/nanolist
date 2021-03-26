package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/ini.v1"
	"io"
	"io/ioutil"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type Config struct {
	CommandAddress string `ini:"command_address"`
	Log            string `ini:"log"`
	Database       string `ini:"database"`
	SMTPHostname   string `ini:"smtp_hostname"`
	SMTPPort       string `ini:"smtp_port"`
	SMTPUsername   string `ini:"smtp_username"`
	SMTPPassword   string `ini:"smtp_password"`
	Lists          map[string]*List
	Debug          bool
	ConfigFile     string
}

type List struct {
	Name            string `ini:"name"`
	Description     string `ini:"description"`
	Id              string
	Address         string   `ini:"address"`
	Hidden          bool     `ini:"hidden"`
	SubscribersOnly bool     `ini:"subscribers_only"`
	Posters         []string `ini:"posters,omitempty"`
	Bcc             []string `ini:"bcc,omitempty"`
}

type Message struct {
	Subject            string
	From               string
	To                 string
	Cc                 string
	Bcc                string
	Date               string
	Id                 string
	InReplyTo          string
	ContentType        string
	ContentDisposition string
	XList              string
	Body               string
	MimeVersion        string
}

var gConfig *Config

// Entry point
func main() {
	gConfig = &Config{}

	flag.BoolVar(&gConfig.Debug, "debug", false, "Don't send emails - print them to stdout instead")
	flag.StringVar(&gConfig.ConfigFile, "config", "", "Load configuration from specified file")
	flag.Parse()

	loadConfig()

	if len(flag.Args()) < 1 {
		fmt.Printf("Error: Command not specified\n")
		os.Exit(1)
	}

	if flag.Arg(0) == "check" {
		if checkConfig() {
			fmt.Printf("Congratulations, nanolist appears to be successfully set up!")
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	requireLog()

	if flag.Arg(0) == "message" {
		msg := &Message{}
		err := msg.FromReader(bufio.NewReader(os.Stdin))
		if err != nil {
			log.Printf("ERROR_PARSING_MESSAGE Error=%q\n", err.Error())
			os.Exit(0)
		}
		log.Printf("MESSAGE_RECEIVED Id=%q From=%q To=%q Cc=%q Bcc=%q Subject=%q\n",
			msg.Id, msg.From, msg.To, msg.Cc, msg.Bcc, msg.Subject)
		handleMessage(msg)
	} else {
		fmt.Printf("Unknown command %s\n", flag.Arg(0))
	}
}

// Figure out if this is a command, or a mailing list post
func handleMessage(msg *Message) {
	if isToCommandAddress(msg) {
		handleCommand(msg)
	} else {
		lists := lookupLists(msg)
		if len(lists) > 0 {
			for _, list := range lists {
				if list.CanPost(msg.From) {
					listMsg := msg.ResendAs(list.Id, list.Address)
					list.Send(listMsg)
					log.Printf("MESSAGE_SENT ListId=%q Id=%q From=%q To=%q Cc=%q Bcc=%q Subject=%q\n",
						list.Id, listMsg.Id, listMsg.From, listMsg.To, listMsg.Cc, listMsg.Bcc, listMsg.Subject)
				} else {
					handleNotAuthorisedToPost(msg)
				}
			}
		} else {
			handleNoDestination(msg)
		}
	}
}

// Handle the command given by the user
func handleCommand(msg *Message) {
	if msg.Subject == "lists" {
		handleShowLists(msg)
	} else if msg.Subject == "help" {
		handleHelp(msg)
	} else if strings.HasPrefix(msg.Subject, "subscribe") {
		handleSubscribe(msg)
	} else if strings.HasPrefix(msg.Subject, "unsubscribe") {
		handleUnsubscribe(msg)
	} else {
		handleUnknownCommand(msg)
	}
}

// Reply to a message that has nowhere to go
func handleNoDestination(msg *Message) {
	reply := msg.Reply()
	reply.From = gConfig.CommandAddress
	reply.Body = "No mailing lists addressed. Your message has not been delivered.\r\n"
	reply.Send([]string{msg.From})
	log.Printf("UNKNOWN_DESTINATION From=%q To=%q Cc=%q Bcc=%q", msg.From, msg.To, msg.Cc, msg.Bcc)
}

// Reply that the user isn't authorised to post to the list
func handleNotAuthorisedToPost(msg *Message) {
	reply := msg.Reply()
	reply.From = gConfig.CommandAddress
	reply.Body = "You are not an approved poster for this mailing list. Your message has not been delivered.\r\n"
	reply.Send([]string{msg.From})
	log.Printf("UNAUTHORISED_POST From=%q To=%q Cc=%q Bcc=%q", msg.From, msg.To, msg.Cc, msg.Bcc)
}

// Reply to an unknown command, giving some help
func handleUnknownCommand(msg *Message) {
	reply := msg.Reply()
	reply.From = gConfig.CommandAddress
	reply.Body = fmt.Sprintf(
		"%s is not a valid command.\r\n\r\n"+
			"Valid commands are:\r\n\r\n"+
			commandInfo(),
		msg.Subject)
	reply.Send([]string{msg.From})
	log.Printf("UNKNOWN_COMMAND From=%q", msg.From)
}

// Reply to a help command with help information
func handleHelp(msg *Message) {
	var body bytes.Buffer
	fmt.Fprintf(&body, commandInfo())
	reply := msg.Reply()
	reply.From = gConfig.CommandAddress
	reply.Body = body.String()
	reply.Send([]string{msg.From})
	log.Printf("HELP_SENT To=%q", reply.To)
}

// Reply to a show mailing lists command with a list of mailing lists
func handleShowLists(msg *Message) {
	var body bytes.Buffer
	fmt.Fprintf(&body, "Available mailing lists:\r\n\r\n")
	for _, list := range gConfig.Lists {
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
		gConfig.CommandAddress)

	reply := msg.Reply()
	reply.From = gConfig.CommandAddress
	reply.Body = body.String()
	reply.Send([]string{msg.From})
	log.Printf("LIST_SENT To=%q", reply.To)
}

// Handle a subscribe command
func handleSubscribe(msg *Message) {
	listId := strings.TrimPrefix(msg.Subject, "subscribe ")
	list := lookupList(listId)

	if list == nil {
		reply := msg.Reply()
		reply.Body = fmt.Sprintf("Unable to subscribe to %s  - it is not a valid mailing list.\r\n", listId)
		reply.Send([]string{msg.From})
		log.Printf("INVALID_SUBSCRIPTION_REQUEST User=%q List=%q\n", msg.From, listId)
		os.Exit(0)
	}

	// Switch to id - in case we were passed address
	listId = list.Id

	if isSubscribed(msg.From, listId) {
		reply := msg.Reply()
		reply.Body = fmt.Sprintf("You are already subscribed to %s\r\n", listId)
		reply.Send([]string{msg.From})
		log.Printf("DUPLICATE_SUBSCRIPTION_REQUEST User=%q List=%q\n", msg.From, listId)
		os.Exit(0)
	}

	addSubscription(msg.From, listId)
	reply := msg.Reply()
	reply.Body = fmt.Sprintf("You are now subscribed to %s\r\n", listId)
	reply.Send([]string{msg.From})
}

// Handle an unsubscribe command
func handleUnsubscribe(msg *Message) {
	listId := strings.TrimPrefix(msg.Subject, "unsubscribe ")
	list := lookupList(listId)

	if list == nil {
		reply := msg.Reply()
		reply.Body = fmt.Sprintf("Unable to unsubscribe from %s  - it is not a valid mailing list.\r\n", listId)
		reply.Send([]string{msg.From})
		log.Printf("INVALID_UNSUBSCRIPTION_REQUEST User=%q List=%q\n", msg.From, listId)
		os.Exit(0)
	}

	// Switch to id - in case we were passed address
	listId = list.Id

	if !isSubscribed(msg.From, listId) {
		reply := msg.Reply()
		reply.Body = fmt.Sprintf("You aren't subscribed to %s\r\n", listId)
		reply.Send([]string{msg.From})
		log.Printf("DUPLICATE_UNSUBSCRIPTION_REQUEST User=%q List=%q\n", msg.From, listId)
		os.Exit(0)
	}

	removeSubscription(msg.From, listId)
	reply := msg.Reply()
	reply.Body = fmt.Sprintf("You are now unsubscribed from %s\r\n", listId)
	reply.Send([]string{msg.From})
}

// MESSAGE LOGIC //////////////////////////////////////////////////////////////

// Read a message from the given io.Reader
func (msg *Message) FromReader(stream io.Reader) error {
	inMessage, err := mail.ReadMessage(stream)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(inMessage.Body)
	if err != nil {
		return err
	}

	msg.MimeVersion = inMessage.Header.Get("MIME-Version")
	msg.Subject = inMessage.Header.Get("Subject")
	msg.From = inMessage.Header.Get("From")
	msg.Id = inMessage.Header.Get("Message-ID")
	msg.InReplyTo = inMessage.Header.Get("In-Reply-To")
	msg.Body = string(body[:])
	msg.To = inMessage.Header.Get("To")
	msg.Cc = inMessage.Header.Get("Cc")
	msg.Bcc = inMessage.Header.Get("Bcc")
	msg.Date = inMessage.Header.Get("Date")
	msg.ContentType = inMessage.Header.Get("Content-Type")
	msg.ContentDisposition = inMessage.Header.Get("Content-Disposition")

	return nil
}

// Create a new message that replies to this message
func (msg *Message) Reply() *Message {
	reply := &Message{}
	reply.Subject = "Re: " + msg.Subject
	reply.To = msg.From
	reply.InReplyTo = msg.Id
	reply.Date = time.Now().Format("Mon, 2 Jan 2006 15:04:05 -0700")
	return reply
}

// Prepare a copy of the message that we're forwarding to a list
func (msg *Message) ResendAs(listId string, listAddress string) *Message {
	send := &Message{}
	send.Subject = msg.Subject
	send.From = msg.From
	send.To = msg.To
	send.Cc = msg.Cc
	send.Date = msg.Date
	send.Id = msg.Id
	send.InReplyTo = msg.InReplyTo
	send.XList = listId + " <" + listAddress + ">"
	send.Body = msg.Body
	send.ContentType = msg.ContentType
	send.ContentDisposition = msg.ContentDisposition
	send.MimeVersion = msg.MimeVersion

	// If the destination mailing list is in the Bcc field, keep it there
	bccList, err := mail.ParseAddressList(msg.Bcc)
	if err == nil {
		for _, bcc := range bccList {
			if bcc.Address == listAddress {
				send.Bcc = listId + " <" + listAddress + ">"
				break
			}
		}
	}
	return send
}

// Generate a emailable represenation of this message
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
	if len(msg.ContentDisposition) > 0 {
		fmt.Fprintf(&buf, "Content-Disposition: %s\r\n", msg.ContentDisposition)
	}
	if len(msg.MimeVersion) > 0 {
		fmt.Fprintf(&buf, "MIME-Version: %s\r\n", msg.MimeVersion)
	}

	fmt.Fprintf(&buf, "Subject: %s\r\n", msg.Subject)
	fmt.Fprintf(&buf, "\r\n%s", msg.Body)

	return buf.String()
}

func (msg *Message) Send(recipients []string) {

	e, err_parse := mail.ParseAddress(msg.From)
	if err_parse != nil {
		log.Printf("ERROR_PARSING Error=%s\n", err_parse)
		os.Exit(0)
	}

	if gConfig.Debug {
		fmt.Printf("------------------------------------------------------------\n")
		fmt.Printf("SENDING MESSAGE TO:\n")
		for _, r := range recipients {
			fmt.Printf(" - %s\n", r)
		}
		fmt.Printf("MESSAGE:\n")
		fmt.Printf("%s\n", msg.String())
		return
	}

	auth := smtp.PlainAuth("", gConfig.SMTPUsername, gConfig.SMTPPassword, gConfig.SMTPHostname)
	err := smtp.SendMail(gConfig.SMTPHostname+":"+gConfig.SMTPPort, auth, e.Address, recipients, []byte(msg.String()))
	if err != nil {
		log.Printf("ERROR_SENDING Error=%q\n", err.Error())
		os.Exit(0)
	}
}

// MAILING LIST LOGIC /////////////////////////////////////////////////////////

// Check if the user is authorised to post to this mailing list
func (list *List) CanPost(from string) bool {

	// Is this list restricted to subscribers only?
	if list.SubscribersOnly && !isSubscribed(from, list.Id) {
		return false
	}

	// Is there a whitelist of approved posters?
	if len(list.Posters) > 0 {
		for _, poster := range list.Posters {
			if from == poster {
				return true
			}
		}
		return false
	}

	return true
}

// Send a message to the mailing list
func (list *List) Send(msg *Message) {
	recipients := fetchSubscribers(list.Id)
	for _, bcc := range list.Bcc {
		recipients = append(recipients, bcc)
	}
	msg.Send(recipients)
}

// DATABASE LOGIC /////////////////////////////////////////////////////////////

// Open the database
func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", gConfig.Database)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS "subscriptions" (
		"list" TEXT,
		"user" TEXT
	);
	`)

	return db, err
}

// Open the database or fail immediately
func requireDB() *sql.DB {
	db, err := openDB()
	if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(1)
	}
	return db
}

// Fetch list of subscribers to a mailing list from database
func fetchSubscribers(listId string) []string {
	db := requireDB()
	rows, err := db.Query("SELECT user FROM subscriptions WHERE list=?", listId)

	if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}

	listIds := []string{}
	defer rows.Close()
	for rows.Next() {
		var user string
		rows.Scan(&user)
		listIds = append(listIds, user)
	}

	return listIds
}

// Check if a user is subscribed to a mailing list
func isSubscribed(user string, list string) bool {
	addressObj, err := mail.ParseAddress(user)
	if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}
	db := requireDB()

	exists := false
	err = db.QueryRow("SELECT 1 FROM subscriptions WHERE user=? AND list=?", addressObj.Address, list).Scan(&exists)

	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}

	return true
}

// Add a subscription to the subscription database
func addSubscription(user string, list string) {
	addressObj, err := mail.ParseAddress(user)
	if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}

	db := requireDB()
	_, err = db.Exec("INSERT INTO subscriptions (user,list) VALUES(?,?)", addressObj.Address, list)
	if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}
	log.Printf("SUBSCRIPTION_ADDED User=%q List=%q\n", user, list)
}

// Remove a subscription from the subscription database
func removeSubscription(user string, list string) {
	addressObj, err := mail.ParseAddress(user)
	if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}

	db := requireDB()
	_, err = db.Exec("DELETE FROM subscriptions WHERE user=? AND list=?", addressObj.Address, list)
	if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}
	log.Printf("SUBSCRIPTION_REMOVED User=%q List=%q\n", user, list)
}

// Remove all subscriptions from a given mailing list
func clearSubscriptions(list string) {
	db := requireDB()
	_, err := db.Exec("DELETE FROM subscriptions WHERE AND list=?", list)
	if err != nil {
		log.Printf("DATABASE_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}
}

// HELPER FUNCTIONS ///////////////////////////////////////////////////////////

// Open the log file for logging
func openLog() error {
	logFile, err := os.OpenFile(gConfig.Log, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	out := io.MultiWriter(logFile, os.Stderr)
	log.SetOutput(out)
	return nil
}

// Open the log, or fail immediately
func requireLog() {
	err := openLog()
	if err != nil {
		log.Printf("LOG_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}
}

// Load gConfig from the on-disk config file
func loadConfig() {
	var (
		err error
		cfg *ini.File
	)

	if len(gConfig.ConfigFile) > 0 {
		cfg, err = ini.Load(gConfig.ConfigFile)
	} else {
		cfg, err = ini.LooseLoad("nanolist.ini", "/usr/local/etc/nanolist.ini", "/etc/nanolist.ini")
	}

	if err != nil {
		log.Printf("CONFIG_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}

	err = cfg.Section("").MapTo(gConfig)
	if err != nil {
		log.Printf("CONFIG_ERROR Error=%q\n", err.Error())
		os.Exit(0)
	}

	gConfig.Lists = make(map[string]*List)

	for _, section := range cfg.ChildSections("list") {
		list := &List{}
		err = section.MapTo(list)
		if err != nil {
			log.Printf("CONFIG_ERROR Error=%q\n", err.Error())
			os.Exit(0)
		}
		list.Id = strings.TrimPrefix(section.Name(), "list.")
		gConfig.Lists[list.Address] = list
	}
}

// Retrieve a list of mailing lists that are recipients of the given message
func lookupLists(msg *Message) []*List {
	lists := []*List{}

	toList, err := mail.ParseAddressList(msg.To)
	if err == nil {
		for _, to := range toList {
			list := lookupList(to.Address)
			if list != nil {
				lists = append(lists, list)
			}
		}
	}

	ccList, err := mail.ParseAddressList(msg.Cc)
	if err == nil {
		for _, cc := range ccList {
			list := lookupList(cc.Address)
			if list != nil {
				lists = append(lists, list)
			}
		}
	}

	bccList, err := mail.ParseAddressList(msg.Bcc)
	if err == nil {
		for _, bcc := range bccList {
			list := lookupList(bcc.Address)
			if list != nil {
				lists = append(lists, list)
			}
		}
	}

	return lists
}

// Look up a mailing list by id or address
func lookupList(listKey string) *List {
	for _, list := range gConfig.Lists {
		if listKey == list.Id || listKey == list.Address {
			return list
		}
	}
	return nil
}

// Is the message bound for our command address?
func isToCommandAddress(msg *Message) bool {
	toList, err := mail.ParseAddressList(msg.To)
	if err == nil {
		for _, to := range toList {
			if to.Address == gConfig.CommandAddress {
				return true
			}
		}
	}

	ccList, err := mail.ParseAddressList(msg.Cc)
	if err == nil {
		for _, cc := range ccList {
			if cc.Address == gConfig.CommandAddress {
				return true
			}
		}
	}

	bccList, err := mail.ParseAddressList(msg.Bcc)
	if err == nil {
		for _, bcc := range bccList {
			if bcc.Address == gConfig.CommandAddress {
				return true
			}
		}
	}

	return false
}

// Generate an email-able list of commands
func commandInfo() string {
	return fmt.Sprintf("    help\r\n"+
		"      Information about valid commands\r\n"+
		"\r\n"+
		"    list\r\n"+
		"      Retrieve a list of available mailing lists\r\n"+
		"\r\n"+
		"    subscribe <list-id>\r\n"+
		"      Subscribe to <list-id>\r\n"+
		"\r\n"+
		"    unsubscribe <list-id>\r\n"+
		"      Unsubscribe from <list-id>\r\n"+
		"\r\n"+
		"To send a command, email %s with the command as the subject.\r\n",
		gConfig.CommandAddress)
}

// Check for a valid configuration
func checkConfig() bool {
	_, err := openDB()
	if err != nil {
		fmt.Printf("There's a problem with the database: %s\n", err.Error())
		return false
	}

	err = openLog()
	if err != nil {
		fmt.Printf("There's a problem with the log: %s\n", err.Error())
		return false
	}

	client, err := smtp.Dial(gConfig.SMTPHostname + ":" + gConfig.SMTPPort)
	if err != nil {
		fmt.Printf("There's a problem connecting to your SMTP server: %s\n", err.Error())
		return false
	}

	auth := smtp.PlainAuth("", gConfig.SMTPUsername, gConfig.SMTPPassword, gConfig.SMTPHostname)
	err = client.Auth(auth)
	if err != nil {
		fmt.Printf("There's a problem authenticating with your SMTP server: %s\n", err.Error())
		return false
	}

	return true
}
