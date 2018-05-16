package main

import (
	"errors"
	"gopkg.in/ini.v1"
	"io"
	"io/ioutil"
	"strings"
)

type INIListLoader struct {
	Reader io.Reader
}

func (ll *INIListLoader) LoadLists(lm ListManager) error {
	data, err := ioutil.ReadAll(ll.Reader)
	if err != nil {
		return err
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return err
	}

	for _, section := range cfg.ChildSections("list") {
		list := &List{}

		if !section.HasKey("address") {
			return errors.New(section.Name() + " has no address")
		}

		list.Name = section.Key("name").String()
		list.Description = section.Key("description").String()
		list.Id = strings.TrimPrefix(section.Name(), "list.")
		list.Address = section.Key("address").String()
		list.Hidden = section.Key("hidden").MustBool(false)
		list.SubscribersOnly = section.Key("subscribers_only").MustBool(false)
		list.Posters = section.Key("posters").Strings(",")
		list.Bcc = section.Key("bcc").Strings(",")

		lm.Add(list)
	}

	return nil
}
