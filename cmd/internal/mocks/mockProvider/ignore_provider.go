// Code generated by MockGen. DO NOT EDIT.
// Source: client/provider/ignore_provider.go

// Package mockProvider is a generated GoMock package.
package mockProvider

import (
	reflect "reflect"

	stevedore "github.com/gojek/stevedore/pkg/stevedore"
	gomock "github.com/golang/mock/gomock"
)

// MockIgnoreProvider is a mock of IgnoreProvider interface.
type MockIgnoreProvider struct {
	ctrl     *gomock.Controller
	recorder *MockIgnoreProviderMockRecorder
}

// MockIgnoreProviderMockRecorder is the mock recorder for MockIgnoreProvider.
type MockIgnoreProviderMockRecorder struct {
	mock *MockIgnoreProvider
}

// NewMockIgnoreProvider creates a new mock instance.
func NewMockIgnoreProvider(ctrl *gomock.Controller) *MockIgnoreProvider {
	mock := &MockIgnoreProvider{ctrl: ctrl}
	mock.recorder = &MockIgnoreProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIgnoreProvider) EXPECT() *MockIgnoreProviderMockRecorder {
	return m.recorder
}

// Files mocks base method.
func (m *MockIgnoreProvider) Files() ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Files")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Files indicates an expected call of Files.
func (mr *MockIgnoreProviderMockRecorder) Files() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Files", reflect.TypeOf((*MockIgnoreProvider)(nil).Files))
}

// Ignores mocks base method.
func (m *MockIgnoreProvider) Ignores() (stevedore.Ignores, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ignores")
	ret0, _ := ret[0].(stevedore.Ignores)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ignores indicates an expected call of Ignores.
func (mr *MockIgnoreProviderMockRecorder) Ignores() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ignores", reflect.TypeOf((*MockIgnoreProvider)(nil).Ignores))
}
