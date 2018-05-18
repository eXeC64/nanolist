package main

import (
	"strings"
	"testing"
)

func TestLoadingLists(t *testing.T) {

	lm := &MemoryListManager{}

	ini := strings.NewReader(`
[list.golang]
# Address this list should receieve mail on
address = golang@example.com
# Information to show in the list of mailing lists
name = "Go programming"
description = "General discussion of Go programming"
# bcc all posts to the listed addresses for archival
bcc = archive@example.com, datahoarder@example.com

[list.announcements]
address = announce@example.com
name = "Announcements"
description = "Important announcements"
# List of email addresses that are permitted to post to this list
posters = admin@example.com, moderator@example.com

[list.fight-club]
address = robertpaulson99@example.com
# Don't tell users this list exists
hidden = true
# Only let subscribed users post to this list
subscribers_only = true`)

	ll := &INIListLoader{ini}

	err := ll.LoadLists(lm)
	if err != nil {
		t.Errorf("Loading lists failed with error: %s", err.Error())
		return
	}

	ids, err := lm.FetchListIds()
	if err != nil {
		t.Errorf("Fetching list ids failed with error: %s", err.Error())
		return
	}
	t.Logf("Loaded the following lists: %v\n", ids)

	golang, err := lm.FetchList("golang@example.com")
	if err != nil {
		t.Errorf("Loading golang list failed with error: %s", err.Error())
	} else {
		// Inspect golang list
		if golang.Id != "golang" {
			t.Errorf("Incorrect golang id: %s", golang.Id)
		}
		if golang.Address != "golang@example.com" {
			t.Errorf("Incorrect golang address: %s", golang.Address)
		}
		if golang.Name != "Go programming" {
			t.Errorf("Incorrect golang name: %s", golang.Name)
		}
		if golang.Description != "General discussion of Go programming" {
			t.Errorf("Incorrect golang description: %s", golang.Description)
		}
		if golang.Hidden != false {
			t.Error("golang incorrectly hidden")
		}
		if golang.SubscribersOnly != false {
			t.Error("golang incorrectly subscribers only")
		}
		if len(golang.Bcc) != 2 {
			t.Errorf("Incorrect golang bcc length: %d", len(golang.Bcc))
		} else {
			if golang.Bcc[0] != "archive@example.com" {
				t.Errorf("Incorrect golang first bcc: %s", golang.Bcc[0])
			}
			if golang.Bcc[1] != "datahoarder@example.com" {
				t.Errorf("Incorrect golang second bcc: %s", golang.Bcc[1])
			}
		}
	}

	announcements, err := lm.FetchList("announce@example.com")
	if err != nil {
		t.Errorf("Loading announcements list failed with error: %s", err.Error())
	} else {
		// Inspect announcements list
		if announcements.Id != "announcements" {
			t.Errorf("Incorrect announcements id: %s", announcements.Id)
		}
		if announcements.Address != "announce@example.com" {
			t.Errorf("Incorrect announcements address: %s", announcements.Address)
		}
		if announcements.Name != "Announcements" {
			t.Errorf("Incorrect announcements name: %s", announcements.Name)
		}
		if announcements.Description != "Important announcements" {
			t.Errorf("Incorrect announcements description: %s", announcements.Description)
		}
		if announcements.Hidden != false {
			t.Error("announcements incorrectly hidden")
		}
		if announcements.SubscribersOnly != false {
			t.Error("announcements incorrectly subscribers only")
		}
		if len(announcements.Posters) != 2 {
			t.Errorf("Incorrect announcements bcc length: %d", len(announcements.Posters))
		} else {
			if announcements.Posters[0] != "admin@example.com" {
				t.Errorf("Incorrect announcements first poster: %s", announcements.Posters[0])
			}
			if announcements.Posters[1] != "moderator@example.com" {
				t.Errorf("Incorrect announcements second poster: %s", announcements.Posters[1])
			}
		}
	}

	fightclub, err := lm.FetchList("robertpaulson99@example.com")
	if err != nil {
		t.Errorf("Loading fightclub list failed with error: %s", err.Error())
	} else {
		// Inspect fightclub list
		if fightclub.Id != "fight-club" {
			t.Errorf("Incorrect fight-club id: %s", fightclub.Id)
		}
		if fightclub.Address != "robertpaulson99@example.com" {
			t.Errorf("Incorrect fightclub address: %s", fightclub.Address)
		}
		if fightclub.Hidden != true {
			t.Error("fightclub should be hidden")
		}
		if fightclub.SubscribersOnly != true {
			t.Error("fightclub should be subscribers only")
		}
	}
}
