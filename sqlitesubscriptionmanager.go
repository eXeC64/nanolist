package main

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteSubscriptionManager struct {
	db *sql.DB
}

func (s *SQLiteSubscriptionManager) Open(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS "subscriptions" (
		"list" TEXT,
		"user" TEXT
	);
	`)

	if err == nil {
		s.db = db
	}
	return err
}

func (s *SQLiteSubscriptionManager) Subscribe(email string, list string) error {
	if s.db == nil {
		return errors.New("No database open")
	}

	_, err := s.db.Exec("INSERT INTO subscriptions (user,list) VALUES(?,?)", email, list)
	return err
}

func (s *SQLiteSubscriptionManager) Unsubscribe(email string, list string) error {
	if s.db == nil {
		return errors.New("No database open")
	}

	_, err := s.db.Exec("DELETE FROM subscriptions WHERE user=? AND list=?", email, list)
	return err
}

func (s *SQLiteSubscriptionManager) UnsubscribeAll(list string) error {
	if s.db == nil {
		return errors.New("No database open")
	}

	_, err := s.db.Exec("DELETE FROM subscriptions WHERE list=?", list)
	return err
}

func (s *SQLiteSubscriptionManager) IsSubscribed(email string, list string) (bool, error) {
	if s.db == nil {
		return false, errors.New("No database open")
	}

	exists := false
	err := s.db.QueryRow("SELECT 1 FROM subscriptions WHERE user=? AND list=?", email, list).Scan(&exists)

	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (s *SQLiteSubscriptionManager) FetchSubscribers(list string) ([]string, error) {
	listAddrs := []string{}

	if s.db == nil {
		return listAddrs, errors.New("No database open")
	}

	rows, err := s.db.Query("SELECT user FROM subscriptions WHERE list=? ORDER BY user ASC", list)

	if err != nil {
		return listAddrs, err
	}

	defer rows.Close()
	for rows.Next() {
		var user string
		rows.Scan(&user)
		listAddrs = append(listAddrs, user)
	}

	return listAddrs, nil
}
