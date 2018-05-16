package main

type List struct {
	Name            string
	Description     string
	Id              string
	Address         string
	Hidden          bool
	SubscribersOnly bool
	Posters         []string
	Bcc             []string
}
