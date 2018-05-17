// Code generated by mockery v1.0.0. DO NOT EDIT.
package main

import mock "github.com/stretchr/testify/mock"

// MockSubscriptionManager is an autogenerated mock type for the SubscriptionManager type
type MockSubscriptionManager struct {
	mock.Mock
}

// FetchSubscribers provides a mock function with given fields: list
func (_m *MockSubscriptionManager) FetchSubscribers(list string) ([]string, error) {
	ret := _m.Called(list)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(list)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(list)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsSubscribed provides a mock function with given fields: email, list
func (_m *MockSubscriptionManager) IsSubscribed(email string, list string) (bool, error) {
	ret := _m.Called(email, list)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(email, list)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(email, list)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Subscribe provides a mock function with given fields: email, list
func (_m *MockSubscriptionManager) Subscribe(email string, list string) error {
	ret := _m.Called(email, list)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(email, list)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Unsubscribe provides a mock function with given fields: email, list
func (_m *MockSubscriptionManager) Unsubscribe(email string, list string) error {
	ret := _m.Called(email, list)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(email, list)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnsubscribeAll provides a mock function with given fields: list
func (_m *MockSubscriptionManager) UnsubscribeAll(list string) error {
	ret := _m.Called(list)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(list)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}