// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	entity "transmitter-artemis/entity"

	mock "github.com/stretchr/testify/mock"
)

// ClientService is an autogenerated mock type for the ClientService type
type ClientService struct {
	mock.Mock
}

// GetAllClientData provides a mock function with given fields:
func (_m *ClientService) GetAllClientData() ([]entity.ClientData, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for GetAllClientData")
	}

	var r0 []entity.ClientData
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]entity.ClientData, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []entity.ClientData); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]entity.ClientData)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewClientService creates a new instance of ClientService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClientService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ClientService {
	mock := &ClientService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
