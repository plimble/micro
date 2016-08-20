package mocks

import "github.com/stretchr/testify/mock"

import "time"

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *Client) Close() {
	_m.Called()
}

// Publish provides a mock function with given fields: subject, v
func (_m *Client) Publish(subject string, v interface{}) error {
	ret := _m.Called(subject, v)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}) error); ok {
		r0 = rf(subject, v)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Request provides a mock function with given fields: subject, req, res, timeout
func (_m *Client) Request(subject string, req interface{}, res interface{}, timeout time.Duration) error {
	ret := _m.Called(subject, req, res, timeout)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, interface{}, interface{}, time.Duration) error); ok {
		r0 = rf(subject, req, res, timeout)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
