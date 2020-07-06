// Code generated by mockery v1.0.0. DO NOT EDIT.

package utclient

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockAPIClient is an autogenerated mock type for the APIClient type
type MockAPIClient struct {
	mock.Mock
}

// Health provides a mock function with given fields: ctx
func (_m *MockAPIClient) Health(ctx context.Context) (int, error) {
	ret := _m.Called(ctx)

	var r0 int
	if rf, ok := ret.Get(0).(func(context.Context) int); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateVisorUptime provides a mock function with given fields: _a0
func (_m *MockAPIClient) UpdateVisorUptime(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
