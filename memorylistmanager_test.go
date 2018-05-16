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

	outList, err := lm.FetchList("test-list")
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

	err = lm.Remove("test-list")
	if err != nil {
		t.Errorf("Removing failed with error: %s", err.Error())
		return
	}
}

func TestFetchListIds(t *testing.T) {
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

	ids, err := lm.FetchListIds()
	if err != nil {
		t.Errorf("Fetching list ids failed with error: %s", err.Error())
		return
	}

	if len(ids) != 2 {
		t.Errorf("Incorrect number of list ids returned: %d", len(ids))
		return
	}

	if ids[0] != "list-a" {
		t.Errorf("Incorrect first list id: %s", ids[0])
	}

	if ids[1] != "list-b" {
		t.Errorf("Incorrect second list id: %s", ids[1])
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
