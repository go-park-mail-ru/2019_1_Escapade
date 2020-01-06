// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	config "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	database "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"

	mock "github.com/stretchr/testify/mock"
)

// UserCaseI is an autogenerated mock type for the UserCaseI type
type UserCaseI struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *UserCaseI) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Get provides a mock function with given fields:
func (_m *UserCaseI) Get() database.Interface {
	ret := _m.Called()

	var r0 database.Interface
	if rf, ok := ret.Get(0).(func() database.Interface); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(database.Interface)
		}
	}

	return r0
}

// Open provides a mock function with given fields: CDB, db
func (_m *UserCaseI) Open(CDB config.Database, db database.Interface) error {
	ret := _m.Called(CDB, db)

	var r0 error
	if rf, ok := ret.Get(0).(func(config.Database, database.Interface) error); ok {
		r0 = rf(CDB, db)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Use provides a mock function with given fields: db
func (_m *UserCaseI) Use(db database.Interface) error {
	ret := _m.Called(db)

	var r0 error
	if rf, ok := ret.Get(0).(func(database.Interface) error); ok {
		r0 = rf(db)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}