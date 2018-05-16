package main

type ListManager interface {
	Add(list *List) error
	Remove(id string) error
	FetchListIds() ([]string, error)
	FetchList(id string) (*List, error)
}
