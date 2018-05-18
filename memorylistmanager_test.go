package main

import (
	"testing"
)

func TestAddAndFetchList(t *testing.T) {
	lm := &MemoryListManager{}

	list := &List{
		Name:            "Test List",
		Description:     "A test list",
		Id:              "test-list",
		Address:         "test@example.com",
		Hidden:          true,
		SubscribersOnly: true,
		Posters:         []string{"a@example.com", "b@example.com"},
		Bcc:             []string{"c@example.com", "d@example.com"},
	}

	err := lm.Add(list)
	if err != nil {
		t.Errorf("Adding failed with error: %s", err.Error())
		return
	}

	outList, err := lm.FetchList("test@example.com")
	if err != nil {
		t.Errorf("Fetching failed with error: %s", err.Error())
		return
	}

	if list != outList {
		t.Error("Fetched list does not match saved one")
	}
}

func TestFetchNonexistantListFails(t *testing.T) {
	lm := &MemoryListManager{}

	_, err := lm.FetchList("no-such-list")
	if err == nil {
		t.Error("Fetching non-existant list did not return error")
	}
}

func TestRemoveList(t *testing.T) {
	lm := &MemoryListManager{}

	list := &List{
		Name:            "Test List",
		Description:     "A test list",
		Id:              "test-list",
		Address:         "test@example.com",
		Hidden:          true,
		SubscribersOnly: true,
		Posters:         []string{"a@example.com", "b@example.com"},
		Bcc:             []string{"c@example.com", "d@example.com"},
	}

	err := lm.Add(list)
	if err != nil {
		t.Errorf("Adding failed with error: %s", err.Error())
		return
	}

	exists, err := lm.IsValidList(list.Address)
	if err != nil {
		t.Errorf("Checking exists failed with error: %s", err.Error())
		return
	}
	if !exists {
		t.Error("List failed IsValid after addition")
	}

	err = lm.Remove("test@example.com")
	if err != nil {
		t.Errorf("Removing failed with error: %s", err.Error())
		return
	}

	exists, err = lm.IsValidList(list.Address)
	if err != nil {
		t.Errorf("Checking exists failed with error: %s", err.Error())
		return
	}
	if exists {
		t.Error("List passed IsValid after removal")
	}
}

func TestFetchListAddresses(t *testing.T) {
	lm := &MemoryListManager{}

	listA := &List{Id: "list-a", Address: "list-a@example.com"}
	listB := &List{Id: "list-b", Address: "list-b@example.com"}

	err := lm.Add(listA)
	if err != nil {
		t.Errorf("Adding list A failed with error: %s", err.Error())
		return
	}

	err = lm.Add(listB)
	if err != nil {
		t.Errorf("Adding list A failed with error: %s", err.Error())
		return
	}

	addrs, err := lm.FetchListAddresses()
	if err != nil {
		t.Errorf("Fetching list addresses failed with error: %s", err.Error())
		return
	}

	if len(addrs) != 2 {
		t.Errorf("Incorrect number of list addresses returned: %d", len(addrs))
		return
	}

	if addrs[0] != "list-a@example.com" {
		t.Errorf("Incorrect first list address: %s", addrs[0])
	}

	if addrs[1] != "list-b@example.com" {
		t.Errorf("Incorrect second list address: %s", addrs[1])
	}
}
