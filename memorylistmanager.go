package main

import (
	"errors"
	"sort"
)

type MemoryListManager struct {
	lists map[string]*List
}

func (m *MemoryListManager) init() {
	m.lists = make(map[string]*List)
}

func (m *MemoryListManager) Add(list *List) error {
	if m.lists == nil {
		m.init()
	}

	m.lists[list.Address] = list
	return nil
}

func (m *MemoryListManager) Remove(address string) error {
	if m.lists == nil {
		m.init()
	}

	delete(m.lists, address)
	return nil
}

func (m *MemoryListManager) IsValidList(address string) (bool, error) {
	_, exists := m.lists[address]
	return exists, nil
}

func (m *MemoryListManager) FetchListIds() ([]string, error) {
	if m.lists == nil {
		m.init()
	}

	addrs := []string{}

	for address, _ := range m.lists {
		addrs = append(addrs, address)
	}

	sort.Strings(addrs)

	return addrs, nil
}

func (m *MemoryListManager) FetchListAddresses() ([]string, error) {
	if m.lists == nil {
		m.init()
	}

	addresses := []string{}

	for _, list := range m.lists {
		addresses = append(addresses, list.Address)
	}

	sort.Strings(addresses)

	return addresses, nil
}

func (m *MemoryListManager) FetchList(address string) (*List, error) {
	if m.lists == nil {
		m.init()
	}

	list, ok := m.lists[address]
	if !ok {
		return nil, errors.New("Invalid list")
	}

	return list, nil
}
