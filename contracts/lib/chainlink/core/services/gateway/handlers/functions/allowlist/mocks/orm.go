// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	mock "github.com/stretchr/testify/mock"
)

// ORM is an autogenerated mock type for the ORM type
type ORM struct {
	mock.Mock
}

type ORM_Expecter struct {
	mock *mock.Mock
}

func (_m *ORM) EXPECT() *ORM_Expecter {
	return &ORM_Expecter{mock: &_m.Mock}
}

// CreateAllowedSenders provides a mock function with given fields: ctx, allowedSenders
func (_m *ORM) CreateAllowedSenders(ctx context.Context, allowedSenders []common.Address) error {
	ret := _m.Called(ctx, allowedSenders)

	if len(ret) == 0 {
		panic("no return value specified for CreateAllowedSenders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.Address) error); ok {
		r0 = rf(ctx, allowedSenders)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ORM_CreateAllowedSenders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateAllowedSenders'
type ORM_CreateAllowedSenders_Call struct {
	*mock.Call
}

// CreateAllowedSenders is a helper method to define mock.On call
//   - ctx context.Context
//   - allowedSenders []common.Address
func (_e *ORM_Expecter) CreateAllowedSenders(ctx interface{}, allowedSenders interface{}) *ORM_CreateAllowedSenders_Call {
	return &ORM_CreateAllowedSenders_Call{Call: _e.mock.On("CreateAllowedSenders", ctx, allowedSenders)}
}

func (_c *ORM_CreateAllowedSenders_Call) Run(run func(ctx context.Context, allowedSenders []common.Address)) *ORM_CreateAllowedSenders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.Address))
	})
	return _c
}

func (_c *ORM_CreateAllowedSenders_Call) Return(_a0 error) *ORM_CreateAllowedSenders_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ORM_CreateAllowedSenders_Call) RunAndReturn(run func(context.Context, []common.Address) error) *ORM_CreateAllowedSenders_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteAllowedSenders provides a mock function with given fields: ctx, blockedSenders
func (_m *ORM) DeleteAllowedSenders(ctx context.Context, blockedSenders []common.Address) error {
	ret := _m.Called(ctx, blockedSenders)

	if len(ret) == 0 {
		panic("no return value specified for DeleteAllowedSenders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []common.Address) error); ok {
		r0 = rf(ctx, blockedSenders)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ORM_DeleteAllowedSenders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteAllowedSenders'
type ORM_DeleteAllowedSenders_Call struct {
	*mock.Call
}

// DeleteAllowedSenders is a helper method to define mock.On call
//   - ctx context.Context
//   - blockedSenders []common.Address
func (_e *ORM_Expecter) DeleteAllowedSenders(ctx interface{}, blockedSenders interface{}) *ORM_DeleteAllowedSenders_Call {
	return &ORM_DeleteAllowedSenders_Call{Call: _e.mock.On("DeleteAllowedSenders", ctx, blockedSenders)}
}

func (_c *ORM_DeleteAllowedSenders_Call) Run(run func(ctx context.Context, blockedSenders []common.Address)) *ORM_DeleteAllowedSenders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]common.Address))
	})
	return _c
}

func (_c *ORM_DeleteAllowedSenders_Call) Return(_a0 error) *ORM_DeleteAllowedSenders_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ORM_DeleteAllowedSenders_Call) RunAndReturn(run func(context.Context, []common.Address) error) *ORM_DeleteAllowedSenders_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllowedSenders provides a mock function with given fields: ctx, offset, limit
func (_m *ORM) GetAllowedSenders(ctx context.Context, offset uint, limit uint) ([]common.Address, error) {
	ret := _m.Called(ctx, offset, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetAllowedSenders")
	}

	var r0 []common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) ([]common.Address, error)); ok {
		return rf(ctx, offset, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, uint) []common.Address); ok {
		r0 = rf(ctx, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, uint) error); ok {
		r1 = rf(ctx, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ORM_GetAllowedSenders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllowedSenders'
type ORM_GetAllowedSenders_Call struct {
	*mock.Call
}

// GetAllowedSenders is a helper method to define mock.On call
//   - ctx context.Context
//   - offset uint
//   - limit uint
func (_e *ORM_Expecter) GetAllowedSenders(ctx interface{}, offset interface{}, limit interface{}) *ORM_GetAllowedSenders_Call {
	return &ORM_GetAllowedSenders_Call{Call: _e.mock.On("GetAllowedSenders", ctx, offset, limit)}
}

func (_c *ORM_GetAllowedSenders_Call) Run(run func(ctx context.Context, offset uint, limit uint)) *ORM_GetAllowedSenders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(uint))
	})
	return _c
}

func (_c *ORM_GetAllowedSenders_Call) Return(_a0 []common.Address, _a1 error) *ORM_GetAllowedSenders_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ORM_GetAllowedSenders_Call) RunAndReturn(run func(context.Context, uint, uint) ([]common.Address, error)) *ORM_GetAllowedSenders_Call {
	_c.Call.Return(run)
	return _c
}

// PurgeAllowedSenders provides a mock function with given fields: ctx
func (_m *ORM) PurgeAllowedSenders(ctx context.Context) error {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for PurgeAllowedSenders")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ORM_PurgeAllowedSenders_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PurgeAllowedSenders'
type ORM_PurgeAllowedSenders_Call struct {
	*mock.Call
}

// PurgeAllowedSenders is a helper method to define mock.On call
//   - ctx context.Context
func (_e *ORM_Expecter) PurgeAllowedSenders(ctx interface{}) *ORM_PurgeAllowedSenders_Call {
	return &ORM_PurgeAllowedSenders_Call{Call: _e.mock.On("PurgeAllowedSenders", ctx)}
}

func (_c *ORM_PurgeAllowedSenders_Call) Run(run func(ctx context.Context)) *ORM_PurgeAllowedSenders_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *ORM_PurgeAllowedSenders_Call) Return(_a0 error) *ORM_PurgeAllowedSenders_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ORM_PurgeAllowedSenders_Call) RunAndReturn(run func(context.Context) error) *ORM_PurgeAllowedSenders_Call {
	_c.Call.Return(run)
	return _c
}

// NewORM creates a new instance of ORM. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewORM(t interface {
	mock.TestingT
	Cleanup(func())
}) *ORM {
	mock := &ORM{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
