// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mustafa-bugra-yildiz/uphitme/repos/task (interfaces: Repo)
//
// Generated by this command:
//
//	mockgen -typed -package task -destination mock.go . Repo
//

// Package task is a generated GoMock package.
package task

import (
	context "context"
	url "net/url"
	reflect "reflect"

	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
	isgomock struct{}
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// Count mocks base method.
func (m *MockRepo) Count(ctx context.Context) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Count", ctx)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Count indicates an expected call of Count.
func (mr *MockRepoMockRecorder) Count(ctx any) *MockRepoCountCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Count", reflect.TypeOf((*MockRepo)(nil).Count), ctx)
	return &MockRepoCountCall{Call: call}
}

// MockRepoCountCall wrap *gomock.Call
type MockRepoCountCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoCountCall) Return(arg0 int, arg1 error) *MockRepoCountCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoCountCall) Do(f func(context.Context) (int, error)) *MockRepoCountCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoCountCall) DoAndReturn(f func(context.Context) (int, error)) *MockRepoCountCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Create mocks base method.
func (m *MockRepo) Create(ctx context.Context, title string, target url.URL, payload Payload) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, title, target, payload)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepoMockRecorder) Create(ctx, title, target, payload any) *MockRepoCreateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepo)(nil).Create), ctx, title, target, payload)
	return &MockRepoCreateCall{Call: call}
}

// MockRepoCreateCall wrap *gomock.Call
type MockRepoCreateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoCreateCall) Return(arg0 uuid.UUID, arg1 error) *MockRepoCreateCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoCreateCall) Do(f func(context.Context, string, url.URL, Payload) (uuid.UUID, error)) *MockRepoCreateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoCreateCall) DoAndReturn(f func(context.Context, string, url.URL, Payload) (uuid.UUID, error)) *MockRepoCreateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Fail mocks base method.
func (m *MockRepo) Fail(ctx context.Context, taskID uuid.UUID, err error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fail", ctx, taskID, err)
	ret0, _ := ret[0].(error)
	return ret0
}

// Fail indicates an expected call of Fail.
func (mr *MockRepoMockRecorder) Fail(ctx, taskID, err any) *MockRepoFailCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fail", reflect.TypeOf((*MockRepo)(nil).Fail), ctx, taskID, err)
	return &MockRepoFailCall{Call: call}
}

// MockRepoFailCall wrap *gomock.Call
type MockRepoFailCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoFailCall) Return(arg0 error) *MockRepoFailCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoFailCall) Do(f func(context.Context, uuid.UUID, error) error) *MockRepoFailCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoFailCall) DoAndReturn(f func(context.Context, uuid.UUID, error) error) *MockRepoFailCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Get mocks base method.
func (m *MockRepo) Get(ctx context.Context, taskID uuid.UUID) (*Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, taskID)
	ret0, _ := ret[0].(*Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepoMockRecorder) Get(ctx, taskID any) *MockRepoGetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepo)(nil).Get), ctx, taskID)
	return &MockRepoGetCall{Call: call}
}

// MockRepoGetCall wrap *gomock.Call
type MockRepoGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoGetCall) Return(arg0 *Task, arg1 error) *MockRepoGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoGetCall) Do(f func(context.Context, uuid.UUID) (*Task, error)) *MockRepoGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoGetCall) DoAndReturn(f func(context.Context, uuid.UUID) (*Task, error)) *MockRepoGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// List mocks base method.
func (m *MockRepo) List(ctx context.Context, page, pageSize int) ([]Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, page, pageSize)
	ret0, _ := ret[0].([]Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockRepoMockRecorder) List(ctx, page, pageSize any) *MockRepoListCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockRepo)(nil).List), ctx, page, pageSize)
	return &MockRepoListCall{Call: call}
}

// MockRepoListCall wrap *gomock.Call
type MockRepoListCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoListCall) Return(arg0 []Task, arg1 error) *MockRepoListCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoListCall) Do(f func(context.Context, int, int) ([]Task, error)) *MockRepoListCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoListCall) DoAndReturn(f func(context.Context, int, int) ([]Task, error)) *MockRepoListCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// ListPending mocks base method.
func (m *MockRepo) ListPending(ctx context.Context) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPending", ctx)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPending indicates an expected call of ListPending.
func (mr *MockRepoMockRecorder) ListPending(ctx any) *MockRepoListPendingCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPending", reflect.TypeOf((*MockRepo)(nil).ListPending), ctx)
	return &MockRepoListPendingCall{Call: call}
}

// MockRepoListPendingCall wrap *gomock.Call
type MockRepoListPendingCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoListPendingCall) Return(arg0 []uuid.UUID, arg1 error) *MockRepoListPendingCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoListPendingCall) Do(f func(context.Context) ([]uuid.UUID, error)) *MockRepoListPendingCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoListPendingCall) DoAndReturn(f func(context.Context) ([]uuid.UUID, error)) *MockRepoListPendingCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Succeed mocks base method.
func (m *MockRepo) Succeed(ctx context.Context, taskID uuid.UUID, code int, body Payload) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Succeed", ctx, taskID, code, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// Succeed indicates an expected call of Succeed.
func (mr *MockRepoMockRecorder) Succeed(ctx, taskID, code, body any) *MockRepoSucceedCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Succeed", reflect.TypeOf((*MockRepo)(nil).Succeed), ctx, taskID, code, body)
	return &MockRepoSucceedCall{Call: call}
}

// MockRepoSucceedCall wrap *gomock.Call
type MockRepoSucceedCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoSucceedCall) Return(arg0 error) *MockRepoSucceedCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoSucceedCall) Do(f func(context.Context, uuid.UUID, int, Payload) error) *MockRepoSucceedCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoSucceedCall) DoAndReturn(f func(context.Context, uuid.UUID, int, Payload) error) *MockRepoSucceedCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
