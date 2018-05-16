package main

import (
	"errors"
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

	m.lists[list.Id] = list
	return nil
}

func (m *MemoryListManager) Remove(id string) error {
	if m.lists == nil {
		m.init()
	}

	delete(m.lists, id)
	return nil
}

func (m *MemoryListManager) FetchListIds() ([]string, error) {
	if m.lists == nil {
		m.init()
	}

	ids := []string{}

	for id, _ := range m.lists {
		ids = append(ids, id)
	}

	return ids, nil
}

func (m *MemoryListManager) FetchListAddresses() ([]string, error) {
	if m.lists == nil {
		m.init()
	}

	addresses := []string{}

	for _, list := range m.lists {
		addresses = append(addresses, list.Address)
	}

	return addresses, nil
}

func (m *MemoryListManager) FetchList(id string) (*List, error) {
	if m.lists == nil {
		m.init()
	}

	list, ok := m.lists[id]
	if !ok {
		return nil, errors.New("Invalid list")
	}

	return list, nil
}
