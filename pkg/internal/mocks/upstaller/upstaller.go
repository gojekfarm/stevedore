// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/stevedore/upstaller.go

// Package upstaller is a generated GoMock package.
package upstaller

import (
	context "context"
	helm "github.com/gojek/stevedore/pkg/helm"
	stevedore "github.com/gojek/stevedore/pkg/stevedore"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	sync "sync"
)

// MockUpstaller is a mock of Upstaller interface
type MockUpstaller struct {
	ctrl     *gomock.Controller
	recorder *MockUpstallerMockRecorder
}

// MockUpstallerMockRecorder is the mock recorder for MockUpstaller
type MockUpstallerMockRecorder struct {
	mock *MockUpstaller
}

// NewMockUpstaller creates a new mock instance
func NewMockUpstaller(ctrl *gomock.Controller) *MockUpstaller {
	mock := &MockUpstaller{ctrl: ctrl}
	mock.recorder = &MockUpstallerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUpstaller) EXPECT() *MockUpstallerMockRecorder {
	return m.recorder
}

// Upstall mocks base method
func (m *MockUpstaller) Upstall(ctx context.Context, client helm.Client, releaseSpecification stevedore.ReleaseSpecification, file string, responseCh chan<- stevedore.Response, proceed chan<- bool, wg *sync.WaitGroup, opts stevedore.Opts, helmTimeout int64) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Upstall", ctx, client, releaseSpecification, file, responseCh, proceed, wg, opts, helmTimeout)
}

// Upstall indicates an expected call of Upstall
func (mr *MockUpstallerMockRecorder) Upstall(ctx, client, releaseSpecification, file, responseCh, proceed, wg, opts, helmTimeout interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upstall", reflect.TypeOf((*MockUpstaller)(nil).Upstall), ctx, client, releaseSpecification, file, responseCh, proceed, wg, opts, helmTimeout)
}
