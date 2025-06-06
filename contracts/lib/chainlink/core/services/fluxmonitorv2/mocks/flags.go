// Code generated by mockery v2.53.0. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"

	generated "github.com/smartcontractkit/chainlink-evm/gethwrappers/generated"

	mock "github.com/stretchr/testify/mock"

	types "github.com/ethereum/go-ethereum/core/types"
)

// Flags is an autogenerated mock type for the Flags type
type Flags struct {
	mock.Mock
}

type Flags_Expecter struct {
	mock *mock.Mock
}

func (_m *Flags) EXPECT() *Flags_Expecter {
	return &Flags_Expecter{mock: &_m.Mock}
}

// Address provides a mock function with no fields
func (_m *Flags) Address() common.Address {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Address")
	}

	var r0 common.Address
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	return r0
}

// Flags_Address_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Address'
type Flags_Address_Call struct {
	*mock.Call
}

// Address is a helper method to define mock.On call
func (_e *Flags_Expecter) Address() *Flags_Address_Call {
	return &Flags_Address_Call{Call: _e.mock.On("Address")}
}

func (_c *Flags_Address_Call) Run(run func()) *Flags_Address_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Flags_Address_Call) Return(_a0 common.Address) *Flags_Address_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Flags_Address_Call) RunAndReturn(run func() common.Address) *Flags_Address_Call {
	_c.Call.Return(run)
	return _c
}

// ContractExists provides a mock function with no fields
func (_m *Flags) ContractExists() bool {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ContractExists")
	}

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// Flags_ContractExists_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ContractExists'
type Flags_ContractExists_Call struct {
	*mock.Call
}

// ContractExists is a helper method to define mock.On call
func (_e *Flags_Expecter) ContractExists() *Flags_ContractExists_Call {
	return &Flags_ContractExists_Call{Call: _e.mock.On("ContractExists")}
}

func (_c *Flags_ContractExists_Call) Run(run func()) *Flags_ContractExists_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *Flags_ContractExists_Call) Return(_a0 bool) *Flags_ContractExists_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Flags_ContractExists_Call) RunAndReturn(run func() bool) *Flags_ContractExists_Call {
	_c.Call.Return(run)
	return _c
}

// IsLowered provides a mock function with given fields: contractAddr
func (_m *Flags) IsLowered(contractAddr common.Address) (bool, error) {
	ret := _m.Called(contractAddr)

	if len(ret) == 0 {
		panic("no return value specified for IsLowered")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Address) (bool, error)); ok {
		return rf(contractAddr)
	}
	if rf, ok := ret.Get(0).(func(common.Address) bool); ok {
		r0 = rf(contractAddr)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(contractAddr)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Flags_IsLowered_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'IsLowered'
type Flags_IsLowered_Call struct {
	*mock.Call
}

// IsLowered is a helper method to define mock.On call
//   - contractAddr common.Address
func (_e *Flags_Expecter) IsLowered(contractAddr interface{}) *Flags_IsLowered_Call {
	return &Flags_IsLowered_Call{Call: _e.mock.On("IsLowered", contractAddr)}
}

func (_c *Flags_IsLowered_Call) Run(run func(contractAddr common.Address)) *Flags_IsLowered_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(common.Address))
	})
	return _c
}

func (_c *Flags_IsLowered_Call) Return(_a0 bool, _a1 error) *Flags_IsLowered_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Flags_IsLowered_Call) RunAndReturn(run func(common.Address) (bool, error)) *Flags_IsLowered_Call {
	_c.Call.Return(run)
	return _c
}

// ParseLog provides a mock function with given fields: log
func (_m *Flags) ParseLog(log types.Log) (generated.AbigenLog, error) {
	ret := _m.Called(log)

	if len(ret) == 0 {
		panic("no return value specified for ParseLog")
	}

	var r0 generated.AbigenLog
	var r1 error
	if rf, ok := ret.Get(0).(func(types.Log) (generated.AbigenLog, error)); ok {
		return rf(log)
	}
	if rf, ok := ret.Get(0).(func(types.Log) generated.AbigenLog); ok {
		r0 = rf(log)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(generated.AbigenLog)
		}
	}

	if rf, ok := ret.Get(1).(func(types.Log) error); ok {
		r1 = rf(log)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Flags_ParseLog_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ParseLog'
type Flags_ParseLog_Call struct {
	*mock.Call
}

// ParseLog is a helper method to define mock.On call
//   - log types.Log
func (_e *Flags_Expecter) ParseLog(log interface{}) *Flags_ParseLog_Call {
	return &Flags_ParseLog_Call{Call: _e.mock.On("ParseLog", log)}
}

func (_c *Flags_ParseLog_Call) Run(run func(log types.Log)) *Flags_ParseLog_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(types.Log))
	})
	return _c
}

func (_c *Flags_ParseLog_Call) Return(_a0 generated.AbigenLog, _a1 error) *Flags_ParseLog_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Flags_ParseLog_Call) RunAndReturn(run func(types.Log) (generated.AbigenLog, error)) *Flags_ParseLog_Call {
	_c.Call.Return(run)
	return _c
}

// NewFlags creates a new instance of Flags. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewFlags(t interface {
	mock.TestingT
	Cleanup(func())
}) *Flags {
	mock := &Flags{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
