// Code generated by mockery v2.11.0. DO NOT EDIT.

package mocker

import (
	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// Generator is an autogenerated mock type for the Generator type
type Generator struct {
	mock.Mock
}

// Generate provides a mock function with given fields:
func (_m *Generator) Generate() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

// NewGenerator creates a new instance of Generator. It also registers a cleanup function to assert the mocks expectations.
func NewGenerator(t testing.TB) *Generator {
	mock := &Generator{}

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
