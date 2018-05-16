package main

type SubscriptionManager interface {
	Subscribe(email string, list string) error
	Unsubscribe(email string, list string) error
	UnsubscribeAll(list string) error
	IsSubscribed(email string, list string) (bool, error)
	FetchSubscribers(list string) ([]string, error)
}
