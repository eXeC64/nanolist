package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

func createTestSQLiteSubscriptionManager() (SubscriptionManager, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	sm := &SQLiteSubscriptionManager{}
	err = sm.Open(db)
	if err != nil {
		return nil, err
	}

	return sm, nil
}

func testSubscribeUnsubscribe(t *testing.T) {
	sm, err := createTestSQLiteSubscriptionManager()
	if err != nil {
		t.Errorf("Failed to create test SQLiteSubscriptionManager: %q", err)
		return
	}
	checkSubscribeUnsubscribe(t, sm)
}

func testUnsubscribeAll(t *testing.T) {
	sm, err := createTestSQLiteSubscriptionManager()
	if err != nil {
		t.Errorf("Failed to create test SQLiteSubscriptionManager: %q", err)
		return
	}
	checkUnsubscribeAll(t, sm)
}

func testFetchSubscribers(t *testing.T) {
	sm, err := createTestSQLiteSubscriptionManager()
	if err != nil {
		t.Errorf("Failed to create test SQLiteSubscriptionManager: %q", err)
		return
	}
	checkFetchSubscribers(t, sm)
}
