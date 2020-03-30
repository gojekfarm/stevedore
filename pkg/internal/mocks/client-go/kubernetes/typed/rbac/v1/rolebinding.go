// Code generated by MockGen. DO NOT EDIT.
// Source: k8s.io/client-go/kubernetes/typed/rbac/v1 (interfaces: RoleBindingInterface)

// Package mock_v1 is a generated GoMock package.
package mock_v1

import (
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/rbac/v1"
	v10 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	reflect "reflect"
)

// MockRoleBindingInterface is a mock of RoleBindingInterface interface
type MockRoleBindingInterface struct {
	ctrl     *gomock.Controller
	recorder *MockRoleBindingInterfaceMockRecorder
}

// MockRoleBindingInterfaceMockRecorder is the mock recorder for MockRoleBindingInterface
type MockRoleBindingInterfaceMockRecorder struct {
	mock *MockRoleBindingInterface
}

// NewMockRoleBindingInterface creates a new mock instance
func NewMockRoleBindingInterface(ctrl *gomock.Controller) *MockRoleBindingInterface {
	mock := &MockRoleBindingInterface{ctrl: ctrl}
	mock.recorder = &MockRoleBindingInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRoleBindingInterface) EXPECT() *MockRoleBindingInterfaceMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockRoleBindingInterface) Create(arg0 *v1.RoleBinding) (*v1.RoleBinding, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0)
	ret0, _ := ret[0].(*v1.RoleBinding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockRoleBindingInterfaceMockRecorder) Create(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRoleBindingInterface)(nil).Create), arg0)
}

// Delete mocks base method
func (m *MockRoleBindingInterface) Delete(arg0 string, arg1 *v10.DeleteOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRoleBindingInterfaceMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRoleBindingInterface)(nil).Delete), arg0, arg1)
}

// DeleteCollection mocks base method
func (m *MockRoleBindingInterface) DeleteCollection(arg0 *v10.DeleteOptions, arg1 v10.ListOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCollection", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCollection indicates an expected call of DeleteCollection
func (mr *MockRoleBindingInterfaceMockRecorder) DeleteCollection(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCollection", reflect.TypeOf((*MockRoleBindingInterface)(nil).DeleteCollection), arg0, arg1)
}

// Get mocks base method
func (m *MockRoleBindingInterface) Get(arg0 string, arg1 v10.GetOptions) (*v1.RoleBinding, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*v1.RoleBinding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockRoleBindingInterfaceMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRoleBindingInterface)(nil).Get), arg0, arg1)
}

// List mocks base method
func (m *MockRoleBindingInterface) List(arg0 v10.ListOptions) (*v1.RoleBindingList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].(*v1.RoleBindingList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockRoleBindingInterfaceMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRoleBindingInterface)(nil).List), arg0)
}

// Patch mocks base method
func (m *MockRoleBindingInterface) Patch(arg0 string, arg1 types.PatchType, arg2 []byte, arg3 ...string) (*v1.RoleBinding, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Patch", varargs...)
	ret0, _ := ret[0].(*v1.RoleBinding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Patch indicates an expected call of Patch
func (mr *MockRoleBindingInterfaceMockRecorder) Patch(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockRoleBindingInterface)(nil).Patch), varargs...)
}

// Update mocks base method
func (m *MockRoleBindingInterface) Update(arg0 *v1.RoleBinding) (*v1.RoleBinding, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0)
	ret0, _ := ret[0].(*v1.RoleBinding)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update
func (mr *MockRoleBindingInterfaceMockRecorder) Update(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRoleBindingInterface)(nil).Update), arg0)
}

// Watch mocks base method
func (m *MockRoleBindingInterface) Watch(arg0 v10.ListOptions) (watch.Interface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Watch", arg0)
	ret0, _ := ret[0].(watch.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Watch indicates an expected call of Watch
func (mr *MockRoleBindingInterfaceMockRecorder) Watch(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Watch", reflect.TypeOf((*MockRoleBindingInterface)(nil).Watch), arg0)
}
