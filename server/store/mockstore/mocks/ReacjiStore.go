// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import (
	reacji "github.com/kaakaa/mattermost-plugin-reacji/server/reacji"
	mock "github.com/stretchr/testify/mock"
)

// ReacjiStore is an autogenerated mock type for the ReacjiStore type
type ReacjiStore struct {
	mock.Mock
}

// ForceUpdate provides a mock function with given fields: new
func (_m *ReacjiStore) ForceUpdate(new *reacji.List) error {
	ret := _m.Called(new)

	var r0 error
	if rf, ok := ret.Get(0).(func(*reacji.List) error); ok {
		r0 = rf(new)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields:
func (_m *ReacjiStore) Get() (*reacji.List, error) {
	ret := _m.Called()

	var r0 *reacji.List
	if rf, ok := ret.Get(0).(func() *reacji.List); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*reacji.List)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: prev, new
func (_m *ReacjiStore) Update(prev *reacji.List, new *reacji.List) error {
	ret := _m.Called(prev, new)

	var r0 error
	if rf, ok := ret.Get(0).(func(*reacji.List, *reacji.List) error); ok {
		r0 = rf(prev, new)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}