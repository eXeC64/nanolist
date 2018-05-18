package main

type ListManager interface {
	Add(list *List) error
	Remove(address string) error
	IsValidList(address string) (bool, error)
	FetchListAddresses() ([]string, error)
	FetchList(address string) (*List, error)
}
