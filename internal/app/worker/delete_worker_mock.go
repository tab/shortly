// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/worker/delete_worker.go
//
// Generated by this command:
//
//	mockgen -source=internal/app/worker/delete_worker.go -destination=internal/app/worker/delete_worker_mock.go -package=worker
//

// Package worker is a generated GoMock package.
package worker

import (
	reflect "reflect"
	dto "shortly/internal/app/dto"

	gomock "go.uber.org/mock/gomock"
)

// MockWorker is a mock of Worker interface.
type MockWorker struct {
	ctrl     *gomock.Controller
	recorder *MockWorkerMockRecorder
	isgomock struct{}
}

// MockWorkerMockRecorder is the mock recorder for MockWorker.
type MockWorkerMockRecorder struct {
	mock *MockWorker
}

// NewMockWorker creates a new mock instance.
func NewMockWorker(ctrl *gomock.Controller) *MockWorker {
	mock := &MockWorker{ctrl: ctrl}
	mock.recorder = &MockWorkerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWorker) EXPECT() *MockWorkerMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockWorker) Add(req dto.BatchDeleteParams) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Add", req)
}

// Add indicates an expected call of Add.
func (mr *MockWorkerMockRecorder) Add(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockWorker)(nil).Add), req)
}

// Start mocks base method.
func (m *MockWorker) Start() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Start")
}

// Start indicates an expected call of Start.
func (mr *MockWorkerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockWorker)(nil).Start))
}

// Stop mocks base method.
func (m *MockWorker) Stop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Stop")
}

// Stop indicates an expected call of Stop.
func (mr *MockWorkerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockWorker)(nil).Stop))
}