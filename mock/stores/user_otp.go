// Code generated by MockGen. DO NOT EDIT.
// Source: internal/stores/user_otp.go

// Package mock_stores is a generated GoMock package.
package mock_stores

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
	dto "tbox_backend/internal/dto"
)

// MockIUserOtpStore is a mock of IUserOtpStore interface
type MockIUserOtpStore struct {
	ctrl     *gomock.Controller
	recorder *MockIUserOtpStoreMockRecorder
}

// MockIUserOtpStoreMockRecorder is the mock recorder for MockIUserOtpStore
type MockIUserOtpStoreMockRecorder struct {
	mock *MockIUserOtpStore
}

// NewMockIUserOtpStore creates a new mock instance
func NewMockIUserOtpStore(ctrl *gomock.Controller) *MockIUserOtpStore {
	mock := &MockIUserOtpStore{ctrl: ctrl}
	mock.recorder = &MockIUserOtpStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIUserOtpStore) EXPECT() *MockIUserOtpStoreMockRecorder {
	return m.recorder
}

// GetByUserID mocks base method
func (m *MockIUserOtpStore) GetByUserID(userID int) (dto.UserOtp, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserID", userID)
	ret0, _ := ret[0].(dto.UserOtp)
	ret1, _ := ret[1].(bool)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetByUserID indicates an expected call of GetByUserID
func (mr *MockIUserOtpStoreMockRecorder) GetByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserID", reflect.TypeOf((*MockIUserOtpStore)(nil).GetByUserID), userID)
}

// Save mocks base method
func (m *MockIUserOtpStore) Save(userOtp dto.UserOtp) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", userOtp)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (mr *MockIUserOtpStoreMockRecorder) Save(userOtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockIUserOtpStore)(nil).Save), userOtp)
}

// UpdateOtp mocks base method
func (m *MockIUserOtpStore) UpdateOtp(userOtp dto.UserOtp) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOtp", userOtp)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOtp indicates an expected call of UpdateOtp
func (mr *MockIUserOtpStoreMockRecorder) UpdateOtp(userOtp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOtp", reflect.TypeOf((*MockIUserOtpStore)(nil).UpdateOtp), userOtp)
}
