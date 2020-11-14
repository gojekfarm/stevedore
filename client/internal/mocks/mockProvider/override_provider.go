// Code generated by MockGen. DO NOT EDIT.
// Source: client/provider/override_provider.go

// Package mockProvider is a generated GoMock package.
package mockProvider

import (
	reflect "reflect"

	stevedore "github.com/gojek/stevedore/pkg/stevedore"
	gomock "github.com/golang/mock/gomock"
)

// MockOverrideProvider is a mock of OverrideProvider interface.
type MockOverrideProvider struct {
	ctrl     *gomock.Controller
	recorder *MockOverrideProviderMockRecorder
}

// MockOverrideProviderMockRecorder is the mock recorder for MockOverrideProvider.
type MockOverrideProviderMockRecorder struct {
	mock *MockOverrideProvider
}

// NewMockOverrideProvider creates a new mock instance.
func NewMockOverrideProvider(ctrl *gomock.Controller) *MockOverrideProvider {
	mock := &MockOverrideProvider{ctrl: ctrl}
	mock.recorder = &MockOverrideProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOverrideProvider) EXPECT() *MockOverrideProviderMockRecorder {
	return m.recorder
}

// Overrides mocks base method.
func (m *MockOverrideProvider) Overrides() (stevedore.Overrides, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Overrides")
	ret0, _ := ret[0].(stevedore.Overrides)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Overrides indicates an expected call of Overrides.
func (mr *MockOverrideProviderMockRecorder) Overrides() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Overrides", reflect.TypeOf((*MockOverrideProvider)(nil).Overrides))
}
