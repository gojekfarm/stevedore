// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/stevedore/file_utils.go

// Package mocks is a generated GoMock package.
package mocks

import (
	os "os"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	afero "github.com/spf13/afero"
)

// MockFileUtils is a mock of FileUtils interface.
type MockFileUtils struct {
	ctrl     *gomock.Controller
	recorder *MockFileUtilsMockRecorder
}

// MockFileUtilsMockRecorder is the mock recorder for MockFileUtils.
type MockFileUtilsMockRecorder struct {
	mock *MockFileUtils
}

// NewMockFileUtils creates a new mock instance.
func NewMockFileUtils(ctrl *gomock.Controller) *MockFileUtils {
	mock := &MockFileUtils{ctrl: ctrl}
	mock.recorder = &MockFileUtilsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileUtils) EXPECT() *MockFileUtilsMockRecorder {
	return m.recorder
}

// Chmod mocks base method.
func (m *MockFileUtils) Chmod(name string, mode os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chmod", name, mode)
	ret0, _ := ret[0].(error)
	return ret0
}

// Chmod indicates an expected call of Chmod.
func (mr *MockFileUtilsMockRecorder) Chmod(name, mode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chmod", reflect.TypeOf((*MockFileUtils)(nil).Chmod), name, mode)
}

// Chown mocks base method.
func (m *MockFileUtils) Chown(name string, uid, gid int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chown", name, uid, gid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Chown indicates an expected call of Chown.
func (mr *MockFileUtilsMockRecorder) Chown(name, uid, gid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chown", reflect.TypeOf((*MockFileUtils)(nil).Chown), name, uid, gid)
}

// Chtimes mocks base method.
func (m *MockFileUtils) Chtimes(name string, atime, mtime time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Chtimes", name, atime, mtime)
	ret0, _ := ret[0].(error)
	return ret0
}

// Chtimes indicates an expected call of Chtimes.
func (mr *MockFileUtilsMockRecorder) Chtimes(name, atime, mtime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Chtimes", reflect.TypeOf((*MockFileUtils)(nil).Chtimes), name, atime, mtime)
}

// Create mocks base method.
func (m *MockFileUtils) Create(name string) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", name)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockFileUtilsMockRecorder) Create(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockFileUtils)(nil).Create), name)
}

// Mkdir mocks base method.
func (m *MockFileUtils) Mkdir(name string, perm os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mkdir", name, perm)
	ret0, _ := ret[0].(error)
	return ret0
}

// Mkdir indicates an expected call of Mkdir.
func (mr *MockFileUtilsMockRecorder) Mkdir(name, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mkdir", reflect.TypeOf((*MockFileUtils)(nil).Mkdir), name, perm)
}

// MkdirAll mocks base method.
func (m *MockFileUtils) MkdirAll(path string, perm os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MkdirAll", path, perm)
	ret0, _ := ret[0].(error)
	return ret0
}

// MkdirAll indicates an expected call of MkdirAll.
func (mr *MockFileUtilsMockRecorder) MkdirAll(path, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MkdirAll", reflect.TypeOf((*MockFileUtils)(nil).MkdirAll), path, perm)
}

// Name mocks base method.
func (m *MockFileUtils) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockFileUtilsMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockFileUtils)(nil).Name))
}

// Open mocks base method.
func (m *MockFileUtils) Open(name string) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Open", name)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Open indicates an expected call of Open.
func (mr *MockFileUtilsMockRecorder) Open(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockFileUtils)(nil).Open), name)
}

// OpenFile mocks base method.
func (m *MockFileUtils) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenFile", name, flag, perm)
	ret0, _ := ret[0].(afero.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenFile indicates an expected call of OpenFile.
func (mr *MockFileUtilsMockRecorder) OpenFile(name, flag, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenFile", reflect.TypeOf((*MockFileUtils)(nil).OpenFile), name, flag, perm)
}

// Remove mocks base method.
func (m *MockFileUtils) Remove(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockFileUtilsMockRecorder) Remove(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockFileUtils)(nil).Remove), name)
}

// RemoveAll mocks base method.
func (m *MockFileUtils) RemoveAll(path string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveAll", path)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveAll indicates an expected call of RemoveAll.
func (mr *MockFileUtilsMockRecorder) RemoveAll(path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveAll", reflect.TypeOf((*MockFileUtils)(nil).RemoveAll), path)
}

// Rename mocks base method.
func (m *MockFileUtils) Rename(oldname, newname string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rename", oldname, newname)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rename indicates an expected call of Rename.
func (mr *MockFileUtilsMockRecorder) Rename(oldname, newname interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rename", reflect.TypeOf((*MockFileUtils)(nil).Rename), oldname, newname)
}

// Stat mocks base method.
func (m *MockFileUtils) Stat(name string) (os.FileInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stat", name)
	ret0, _ := ret[0].(os.FileInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Stat indicates an expected call of Stat.
func (mr *MockFileUtilsMockRecorder) Stat(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stat", reflect.TypeOf((*MockFileUtils)(nil).Stat), name)
}

// TempDir mocks base method.
func (m *MockFileUtils) TempDir(prefix string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TempDir", prefix)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TempDir indicates an expected call of TempDir.
func (mr *MockFileUtilsMockRecorder) TempDir(prefix interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TempDir", reflect.TypeOf((*MockFileUtils)(nil).TempDir), prefix)
}

// WriteFile mocks base method.
func (m *MockFileUtils) WriteFile(filename string, data []byte, perm os.FileMode) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WriteFile", filename, data, perm)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteFile indicates an expected call of WriteFile.
func (mr *MockFileUtilsMockRecorder) WriteFile(filename, data, perm interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteFile", reflect.TypeOf((*MockFileUtils)(nil).WriteFile), filename, data, perm)
}
