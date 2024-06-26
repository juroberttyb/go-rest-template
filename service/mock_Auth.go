// Code generated by mockery v2.41.0. DO NOT EDIT.

package service

import (
	context "context"

	models "github.com/A-pen-app/kickstart/models"
	mock "github.com/stretchr/testify/mock"
)

// MockAuth is an autogenerated mock type for the Auth type
type MockAuth struct {
	mock.Mock
}

// IssueToken provides a mock function with given fields: ctx, userID, userType, options
func (_m *MockAuth) IssueToken(ctx context.Context, userID string, userType models.UserType, options ...IssueOption) (string, error) {
	_va := make([]interface{}, len(options))
	for _i := range options {
		_va[_i] = options[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, userID, userType)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for IssueToken")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, models.UserType, ...IssueOption) (string, error)); ok {
		return rf(ctx, userID, userType, options...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, models.UserType, ...IssueOption) string); ok {
		r0 = rf(ctx, userID, userType, options...)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, models.UserType, ...IssueOption) error); ok {
		r1 = rf(ctx, userID, userType, options...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ValidateToken provides a mock function with given fields: ctx, token
func (_m *MockAuth) ValidateToken(ctx context.Context, token string) (*models.Claims, error) {
	ret := _m.Called(ctx, token)

	if len(ret) == 0 {
		panic("no return value specified for ValidateToken")
	}

	var r0 *models.Claims
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*models.Claims, error)); ok {
		return rf(ctx, token)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *models.Claims); ok {
		r0 = rf(ctx, token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Claims)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, token)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewMockAuth creates a new instance of MockAuth. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAuth(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAuth {
	mock := &MockAuth{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
