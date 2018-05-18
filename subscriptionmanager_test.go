package main

import (
	"sort"
	"testing"
)

func checkSubscribeUnsubscribe(t *testing.T, sm SubscriptionManager) {
	const (
		addr string = "test@example.com"
		list string = "cool-list@example.com"
	)

	isSubscribed, err := sm.IsSubscribed(addr, list)
	if err != nil {
		t.Errorf("Failed to check subscription: %q", err.Error())
		return
	}

	if isSubscribed {
		t.Error("IsSubscribed was true, expected false")
	}

	err = sm.Subscribe(addr, list)
	if err != nil {
		t.Errorf("Failed to add subscription: %q", err.Error())
		return
	}

	isSubscribed, err = sm.IsSubscribed(addr, list)
	if err != nil {
		t.Errorf("Failed to check subscription: %q", err.Error())
		return
	}

	if !isSubscribed {
		t.Error("IsSubscribed was false, expected true")
	}

	err = sm.Unsubscribe(addr, list)
	if err != nil {
		t.Errorf("Failed to remove subscription: %q", err.Error())
		return
	}

	isSubscribed, err = sm.IsSubscribed(addr, list)
	if err != nil {
		t.Errorf("Failed to check subscription: %q", err.Error())
		return
	}

	if isSubscribed {
		t.Error("IsSubscribed was true, expected false")
	}
}

func checkUnsubscribeAll(t *testing.T, sm SubscriptionManager) {
	const (
		addrA string = "arnold@example.com"
		addrB string = "brook@example.com"
		addrC string = "charlie@example.com"
		listA string = "cool-list@example.com"
		listB string = "other-list@example.com"
	)

	subscriptions := []struct {
		addr string
		list string
	}{
		{addrA, listA},
		{addrB, listA},
		{addrC, listA},
		{addrA, listB},
		{addrC, listB},
	}

	for _, sub := range subscriptions {
		err := sm.Subscribe(sub.addr, sub.list)
		if err != nil {
			t.Errorf("Failed to add subscription: %q", err.Error())
			return
		}
	}

	sm.UnsubscribeAll(listA)

	for _, sub := range subscriptions {
		isSubscribed, err := sm.IsSubscribed(sub.addr, sub.list)
		if err != nil {
			t.Errorf("Failed to check subscription: %q", err.Error())
			return
		}

		if sub.list == listA && isSubscribed {
			t.Errorf(sub.addr + " is still subscribed to " + sub.list)
		} else if sub.list == listB && !isSubscribed {
			t.Errorf(sub.addr + " is not subscribed to " + sub.list)
		}
	}
}

func checkFetchSubscribers(t *testing.T, sm SubscriptionManager) {
	const (
		addrA string = "arnold@example.com"
		addrB string = "brook@example.com"
		addrC string = "charlie@example.com"
		listA string = "cool-list@example.com"
		listB string = "other-list@example.com"
	)

	subscriptions := []struct {
		addr string
		list string
	}{
		{addrA, listA},
		{addrB, listA},
		{addrC, listA},
		{addrA, listB},
		{addrC, listB},
	}

	for _, sub := range subscriptions {
		err := sm.Subscribe(sub.addr, sub.list)
		if err != nil {
			t.Errorf("Failed to add subscription: %q", err.Error())
			return
		}
	}

	results, err := sm.FetchSubscribers(listB)
	if err != nil {
		t.Errorf("Failed to fetch subscriptions: %q", err.Error())
		return
	}

	if len(results) != 2 {
		t.Errorf("Incorrect number of results returned. Expected: 2 Actual: %d", len(results))
	} else {
		sort.Strings(results)
		if results[0] != addrA {
			t.Errorf("Incorrect subscriber. Expected: %q Actual: %q", addrA, results[0])
		}
		if results[1] != addrA {
			t.Errorf("Incorrect subscriber. Expected: %q Actual: %q", addrB, results[1])
		}
	}
}
