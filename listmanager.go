package main

type ListManager interface {
	Add(list *List) error
	Remove(id string) error
	IsValidList(id string) (bool, error)
	FetchListIds() ([]string, error)
	FetchListAddresses() ([]string, error)
	FetchList(id string) (*List, error)
}
