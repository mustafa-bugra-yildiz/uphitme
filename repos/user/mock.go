// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mustafa-bugra-yildiz/uphitme/repos/user (interfaces: Repo)
//
// Generated by this command:
//
//	mockgen -typed -package user -destination mock.go . Repo
//

// Package user is a generated GoMock package.
package user

import (
	context "context"
	reflect "reflect"

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

// Create mocks base method.
func (m *MockRepo) Create(ctx context.Context, fullName, email, hashedPassword string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, fullName, email, hashedPassword)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockRepoMockRecorder) Create(ctx, fullName, email, hashedPassword any) *MockRepoCreateCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepo)(nil).Create), ctx, fullName, email, hashedPassword)
	return &MockRepoCreateCall{Call: call}
}

// MockRepoCreateCall wrap *gomock.Call
type MockRepoCreateCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoCreateCall) Return(arg0 error) *MockRepoCreateCall {
	c.Call = c.Call.Return(arg0)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoCreateCall) Do(f func(context.Context, string, string, string) error) *MockRepoCreateCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoCreateCall) DoAndReturn(f func(context.Context, string, string, string) error) *MockRepoCreateCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}

// Get mocks base method.
func (m *MockRepo) Get(ctx context.Context, email string) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, email)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepoMockRecorder) Get(ctx, email any) *MockRepoGetCall {
	mr.mock.ctrl.T.Helper()
	call := mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepo)(nil).Get), ctx, email)
	return &MockRepoGetCall{Call: call}
}

// MockRepoGetCall wrap *gomock.Call
type MockRepoGetCall struct {
	*gomock.Call
}

// Return rewrite *gomock.Call.Return
func (c *MockRepoGetCall) Return(arg0 *User, arg1 error) *MockRepoGetCall {
	c.Call = c.Call.Return(arg0, arg1)
	return c
}

// Do rewrite *gomock.Call.Do
func (c *MockRepoGetCall) Do(f func(context.Context, string) (*User, error)) *MockRepoGetCall {
	c.Call = c.Call.Do(f)
	return c
}

// DoAndReturn rewrite *gomock.Call.DoAndReturn
func (c *MockRepoGetCall) DoAndReturn(f func(context.Context, string) (*User, error)) *MockRepoGetCall {
	c.Call = c.Call.DoAndReturn(f)
	return c
}
