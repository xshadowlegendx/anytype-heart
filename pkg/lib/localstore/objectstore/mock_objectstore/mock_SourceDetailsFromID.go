// Code generated by mockery v2.38.0. DO NOT EDIT.

package mock_objectstore

import (
	mock "github.com/stretchr/testify/mock"

	types "github.com/gogo/protobuf/types"
)

// MockSourceDetailsFromID is an autogenerated mock type for the SourceDetailsFromID type
type MockSourceDetailsFromID struct {
	mock.Mock
}

type MockSourceDetailsFromID_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSourceDetailsFromID) EXPECT() *MockSourceDetailsFromID_Expecter {
	return &MockSourceDetailsFromID_Expecter{mock: &_m.Mock}
}

// DetailsFromIdBasedSource provides a mock function with given fields: id
func (_m *MockSourceDetailsFromID) DetailsFromIdBasedSource(id string) (*types.Struct, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for DetailsFromIdBasedSource")
	}

	var r0 *types.Struct
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*types.Struct, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) *types.Struct); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.Struct)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSourceDetailsFromID_DetailsFromIdBasedSource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DetailsFromIdBasedSource'
type MockSourceDetailsFromID_DetailsFromIdBasedSource_Call struct {
	*mock.Call
}

// DetailsFromIdBasedSource is a helper method to define mock.On call
//   - id string
func (_e *MockSourceDetailsFromID_Expecter) DetailsFromIdBasedSource(id interface{}) *MockSourceDetailsFromID_DetailsFromIdBasedSource_Call {
	return &MockSourceDetailsFromID_DetailsFromIdBasedSource_Call{Call: _e.mock.On("DetailsFromIdBasedSource", id)}
}

func (_c *MockSourceDetailsFromID_DetailsFromIdBasedSource_Call) Run(run func(id string)) *MockSourceDetailsFromID_DetailsFromIdBasedSource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *MockSourceDetailsFromID_DetailsFromIdBasedSource_Call) Return(_a0 *types.Struct, _a1 error) *MockSourceDetailsFromID_DetailsFromIdBasedSource_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSourceDetailsFromID_DetailsFromIdBasedSource_Call) RunAndReturn(run func(string) (*types.Struct, error)) *MockSourceDetailsFromID_DetailsFromIdBasedSource_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSourceDetailsFromID creates a new instance of MockSourceDetailsFromID. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSourceDetailsFromID(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSourceDetailsFromID {
	mock := &MockSourceDetailsFromID{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}