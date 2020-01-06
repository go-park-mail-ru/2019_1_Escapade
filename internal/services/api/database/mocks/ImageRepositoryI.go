// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	database "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"
	mock "github.com/stretchr/testify/mock"
)

// ImageRepositoryI is an autogenerated mock type for the ImageRepositoryI type
type ImageRepositoryI struct {
	mock.Mock
}

// FetchByID provides a mock function with given fields: dbI, id
func (_m *ImageRepositoryI) FetchByID(dbI database.Interface, id int32) (string, error) {
	ret := _m.Called(dbI, id)

	var r0 string
	if rf, ok := ret.Get(0).(func(database.Interface, int32) string); ok {
		r0 = rf(dbI, id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(database.Interface, int32) error); ok {
		r1 = rf(dbI, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FetchByName provides a mock function with given fields: dbI, name
func (_m *ImageRepositoryI) FetchByName(dbI database.Interface, name string) (string, error) {
	ret := _m.Called(dbI, name)

	var r0 string
	if rf, ok := ret.Get(0).(func(database.Interface, string) string); ok {
		r0 = rf(dbI, name)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(database.Interface, string) error); ok {
		r1 = rf(dbI, name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: dbI, filename, userID
func (_m *ImageRepositoryI) Update(dbI database.Interface, filename string, userID int32) error {
	ret := _m.Called(dbI, filename, userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(database.Interface, string, int32) error); ok {
		r0 = rf(dbI, filename, userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}