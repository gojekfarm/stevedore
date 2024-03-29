// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/stevedore/chart_manager.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	chart "helm.sh/helm/v3/pkg/chart"
)

// MockChartManager is a mock of ChartManager interface.
type MockChartManager struct {
	ctrl     *gomock.Controller
	recorder *MockChartManagerMockRecorder
}

// MockChartManagerMockRecorder is the mock recorder for MockChartManager.
type MockChartManagerMockRecorder struct {
	mock *MockChartManager
}

// NewMockChartManager creates a new mock instance.
func NewMockChartManager(ctrl *gomock.Controller) *MockChartManager {
	mock := &MockChartManager{ctrl: ctrl}
	mock.recorder = &MockChartManagerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockChartManager) EXPECT() *MockChartManagerMockRecorder {
	return m.recorder
}

// Archive mocks base method.
func (m *MockChartManager) Archive(ctx context.Context, ch *chart.Chart, outDir string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Archive", ctx, ch, outDir)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Archive indicates an expected call of Archive.
func (mr *MockChartManagerMockRecorder) Archive(ctx, ch, outDir interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Archive", reflect.TypeOf((*MockChartManager)(nil).Archive), ctx, ch, outDir)
}

// Build mocks base method.
func (m *MockChartManager) Build(ctx context.Context, chartPath string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Build", ctx, chartPath)
	ret0, _ := ret[0].(error)
	return ret0
}

// Build indicates an expected call of Build.
func (mr *MockChartManagerMockRecorder) Build(ctx, chartPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockChartManager)(nil).Build), ctx, chartPath)
}

// Load mocks base method.
func (m *MockChartManager) Load(ctx context.Context, chartPath string) (*chart.Chart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Load", ctx, chartPath)
	ret0, _ := ret[0].(*chart.Chart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Load indicates an expected call of Load.
func (mr *MockChartManagerMockRecorder) Load(ctx, chartPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Load", reflect.TypeOf((*MockChartManager)(nil).Load), ctx, chartPath)
}

// UploadChart mocks base method.
func (m *MockChartManager) UploadChart(ctx context.Context, name, url string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadChart", ctx, name, url)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadChart indicates an expected call of UploadChart.
func (mr *MockChartManagerMockRecorder) UploadChart(ctx, name, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadChart", reflect.TypeOf((*MockChartManager)(nil).UploadChart), ctx, name, url)
}
