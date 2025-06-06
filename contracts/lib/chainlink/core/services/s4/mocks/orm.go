// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	context "context"

	big "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	mock "github.com/stretchr/testify/mock"

	s4 "github.com/smartcontractkit/chainlink/v2/core/services/s4"

	time "time"
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

// DeleteExpired provides a mock function with given fields: ctx, limit, utcNow
func (_m *ORM) DeleteExpired(ctx context.Context, limit uint, utcNow time.Time) (int64, error) {
	ret := _m.Called(ctx, limit, utcNow)

	if len(ret) == 0 {
		panic("no return value specified for DeleteExpired")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, time.Time) (int64, error)); ok {
		return rf(ctx, limit, utcNow)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint, time.Time) int64); ok {
		r0 = rf(ctx, limit, utcNow)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint, time.Time) error); ok {
		r1 = rf(ctx, limit, utcNow)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ORM_DeleteExpired_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteExpired'
type ORM_DeleteExpired_Call struct {
	*mock.Call
}

// DeleteExpired is a helper method to define mock.On call
//   - ctx context.Context
//   - limit uint
//   - utcNow time.Time
func (_e *ORM_Expecter) DeleteExpired(ctx interface{}, limit interface{}, utcNow interface{}) *ORM_DeleteExpired_Call {
	return &ORM_DeleteExpired_Call{Call: _e.mock.On("DeleteExpired", ctx, limit, utcNow)}
}

func (_c *ORM_DeleteExpired_Call) Run(run func(ctx context.Context, limit uint, utcNow time.Time)) *ORM_DeleteExpired_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint), args[2].(time.Time))
	})
	return _c
}

func (_c *ORM_DeleteExpired_Call) Return(_a0 int64, _a1 error) *ORM_DeleteExpired_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ORM_DeleteExpired_Call) RunAndReturn(run func(context.Context, uint, time.Time) (int64, error)) *ORM_DeleteExpired_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, address, slotId
func (_m *ORM) Get(ctx context.Context, address *big.Big, slotId uint) (*s4.Row, error) {
	ret := _m.Called(ctx, address, slotId)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 *s4.Row
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *big.Big, uint) (*s4.Row, error)); ok {
		return rf(ctx, address, slotId)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *big.Big, uint) *s4.Row); ok {
		r0 = rf(ctx, address, slotId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*s4.Row)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *big.Big, uint) error); ok {
		r1 = rf(ctx, address, slotId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ORM_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type ORM_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - address *big.Big
//   - slotId uint
func (_e *ORM_Expecter) Get(ctx interface{}, address interface{}, slotId interface{}) *ORM_Get_Call {
	return &ORM_Get_Call{Call: _e.mock.On("Get", ctx, address, slotId)}
}

func (_c *ORM_Get_Call) Run(run func(ctx context.Context, address *big.Big, slotId uint)) *ORM_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*big.Big), args[2].(uint))
	})
	return _c
}

func (_c *ORM_Get_Call) Return(_a0 *s4.Row, _a1 error) *ORM_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ORM_Get_Call) RunAndReturn(run func(context.Context, *big.Big, uint) (*s4.Row, error)) *ORM_Get_Call {
	_c.Call.Return(run)
	return _c
}

// GetSnapshot provides a mock function with given fields: ctx, addressRange
func (_m *ORM) GetSnapshot(ctx context.Context, addressRange *s4.AddressRange) ([]*s4.SnapshotRow, error) {
	ret := _m.Called(ctx, addressRange)

	if len(ret) == 0 {
		panic("no return value specified for GetSnapshot")
	}

	var r0 []*s4.SnapshotRow
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *s4.AddressRange) ([]*s4.SnapshotRow, error)); ok {
		return rf(ctx, addressRange)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *s4.AddressRange) []*s4.SnapshotRow); ok {
		r0 = rf(ctx, addressRange)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*s4.SnapshotRow)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *s4.AddressRange) error); ok {
		r1 = rf(ctx, addressRange)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ORM_GetSnapshot_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetSnapshot'
type ORM_GetSnapshot_Call struct {
	*mock.Call
}

// GetSnapshot is a helper method to define mock.On call
//   - ctx context.Context
//   - addressRange *s4.AddressRange
func (_e *ORM_Expecter) GetSnapshot(ctx interface{}, addressRange interface{}) *ORM_GetSnapshot_Call {
	return &ORM_GetSnapshot_Call{Call: _e.mock.On("GetSnapshot", ctx, addressRange)}
}

func (_c *ORM_GetSnapshot_Call) Run(run func(ctx context.Context, addressRange *s4.AddressRange)) *ORM_GetSnapshot_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*s4.AddressRange))
	})
	return _c
}

func (_c *ORM_GetSnapshot_Call) Return(_a0 []*s4.SnapshotRow, _a1 error) *ORM_GetSnapshot_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ORM_GetSnapshot_Call) RunAndReturn(run func(context.Context, *s4.AddressRange) ([]*s4.SnapshotRow, error)) *ORM_GetSnapshot_Call {
	_c.Call.Return(run)
	return _c
}

// GetUnconfirmedRows provides a mock function with given fields: ctx, limit
func (_m *ORM) GetUnconfirmedRows(ctx context.Context, limit uint) ([]*s4.Row, error) {
	ret := _m.Called(ctx, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetUnconfirmedRows")
	}

	var r0 []*s4.Row
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) ([]*s4.Row, error)); ok {
		return rf(ctx, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) []*s4.Row); ok {
		r0 = rf(ctx, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*s4.Row)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ORM_GetUnconfirmedRows_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUnconfirmedRows'
type ORM_GetUnconfirmedRows_Call struct {
	*mock.Call
}

// GetUnconfirmedRows is a helper method to define mock.On call
//   - ctx context.Context
//   - limit uint
func (_e *ORM_Expecter) GetUnconfirmedRows(ctx interface{}, limit interface{}) *ORM_GetUnconfirmedRows_Call {
	return &ORM_GetUnconfirmedRows_Call{Call: _e.mock.On("GetUnconfirmedRows", ctx, limit)}
}

func (_c *ORM_GetUnconfirmedRows_Call) Run(run func(ctx context.Context, limit uint)) *ORM_GetUnconfirmedRows_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uint))
	})
	return _c
}

func (_c *ORM_GetUnconfirmedRows_Call) Return(_a0 []*s4.Row, _a1 error) *ORM_GetUnconfirmedRows_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ORM_GetUnconfirmedRows_Call) RunAndReturn(run func(context.Context, uint) ([]*s4.Row, error)) *ORM_GetUnconfirmedRows_Call {
	_c.Call.Return(run)
	return _c
}

// Update provides a mock function with given fields: ctx, row
func (_m *ORM) Update(ctx context.Context, row *s4.Row) error {
	ret := _m.Called(ctx, row)

	if len(ret) == 0 {
		panic("no return value specified for Update")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *s4.Row) error); ok {
		r0 = rf(ctx, row)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ORM_Update_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Update'
type ORM_Update_Call struct {
	*mock.Call
}

// Update is a helper method to define mock.On call
//   - ctx context.Context
//   - row *s4.Row
func (_e *ORM_Expecter) Update(ctx interface{}, row interface{}) *ORM_Update_Call {
	return &ORM_Update_Call{Call: _e.mock.On("Update", ctx, row)}
}

func (_c *ORM_Update_Call) Run(run func(ctx context.Context, row *s4.Row)) *ORM_Update_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*s4.Row))
	})
	return _c
}

func (_c *ORM_Update_Call) Return(_a0 error) *ORM_Update_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ORM_Update_Call) RunAndReturn(run func(context.Context, *s4.Row) error) *ORM_Update_Call {
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
