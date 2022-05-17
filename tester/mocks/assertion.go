// Code generated by MockGen. DO NOT EDIT.
// Source: asserter_2.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	tester "github.com/byorty/contractor/tester"
	gomock "github.com/golang/mock/gomock"
)

// MockAssertionFactory is a mock of Assertion2Factory interface.
type MockAssertionFactory struct {
	ctrl     *gomock.Controller
	recorder *MockAssertionFactoryMockRecorder
}

// MockAssertionFactoryMockRecorder is the mock recorder for MockAssertionFactory.
type MockAssertionFactoryMockRecorder struct {
	mock *MockAssertionFactory
}

// NewMockAssertionFactory creates a new mock instance.
func NewMockAssertionFactory(ctrl *gomock.Controller) *MockAssertionFactory {
	mock := &MockAssertionFactory{ctrl: ctrl}
	mock.recorder = &MockAssertionFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAssertionFactory) EXPECT() *MockAssertionFactoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAssertionFactory) Create(name string, definition interface{}) (tester.Asserter2, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name, definition)
	ret0, _ := ret[0].(tester.Asserter2)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockAssertionFactoryMockRecorder) Create(name, definition interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAssertionFactory)(nil).Create), name, definition)
}

// MockAssertion2 is a mock of Asserter2 interface.
type MockAssertion2 struct {
	ctrl     *gomock.Controller
	recorder *MockAssertion2MockRecorder
}

// MockAssertion2MockRecorder is the mock recorder for MockAssertion2.
type MockAssertion2MockRecorder struct {
	mock *MockAssertion2
}

// NewMockAssertion2 creates a new mock instance.
func NewMockAssertion2(ctrl *gomock.Controller) *MockAssertion2 {
	mock := &MockAssertion2{ctrl: ctrl}
	mock.recorder = &MockAssertion2MockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAssertion2) EXPECT() *MockAssertion2MockRecorder {
	return m.recorder
}

// Assert mocks base method.
func (m *MockAssertion2) Assert(data interface{}) tester.AssertionResultList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Assert", data)
	ret0, _ := ret[0].(tester.AssertionResultList)
	return ret0
}

// Assert indicates an expected call of Assert.
func (mr *MockAssertion2MockRecorder) Assert(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Assert", reflect.TypeOf((*MockAssertion2)(nil).Assert), data)
}

// MockAssertionComparator is a mock of AssertionComparator interface.
type MockAssertionComparator struct {
	ctrl     *gomock.Controller
	recorder *MockAssertionComparatorMockRecorder
}

// MockAssertionComparatorMockRecorder is the mock recorder for MockAssertionComparator.
type MockAssertionComparatorMockRecorder struct {
	mock *MockAssertionComparator
}

// NewMockAssertionComparator creates a new mock instance.
func NewMockAssertionComparator(ctrl *gomock.Controller) *MockAssertionComparator {
	mock := &MockAssertionComparator{ctrl: ctrl}
	mock.recorder = &MockAssertionComparatorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAssertionComparator) EXPECT() *MockAssertionComparatorMockRecorder {
	return m.recorder
}

// Compare mocks base method.
func (m *MockAssertionComparator) Compare() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Compare")
}

// Compare indicates an expected call of Compare.
func (mr *MockAssertionComparatorMockRecorder) Compare() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compare", reflect.TypeOf((*MockAssertionComparator)(nil).Compare))
}