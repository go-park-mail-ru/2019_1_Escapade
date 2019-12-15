package server

import (
	"fmt"
	"os"
)

type InputI interface {
	//WithExtraI

	Init()
	Port() int
	CheckBefore() error
	CheckAfter() error

	GetData() InputData
}

type InputData struct {
	MainPort    string
	MainPortInt int
	FieldPath   string
	RoomPath    string
}

// Input data from user(set by funcs where server launchs)
// If you dont set any callback no error will happen
type Input struct {
	//Extra

	Data InputData

	//set callbacks
	CallInit        func()
	CallCheckBefore func() error
	CallCheckAfter  func() error
}

// Init implements InputI interface
func (input *Input) Init() {
	if input.CallInit == nil {
		return
	}
	input.CallInit()
}

// CheckBefore implements InputI interface
func (input *Input) CheckBefore() error {
	if input.CallCheckBefore == nil {
		return nil
	}
	return input.CallCheckBefore()
}

// CheckAfter implements InputI interface
func (input *Input) CheckAfter() error {
	if input.CallCheckAfter == nil {
		return nil
	}
	return input.CallCheckAfter()
}

// Port implements InputI interface
func (input *Input) GetData() InputData {
	return input.Data
}

// CheckBeforeDefault default implementation of CheckBefore
func (input *Input) CheckBeforeDefault(argsNeed int) error {
	num := len(os.Args)
	if num == argsNeed {
		return nil
	}
	return fmt.Errorf("incorrect amount of arguments. Expected:%d. Get:%d", argsNeed, num)
}

func (input *Input) Port() int {
	return input.Data.MainPortInt
}

// CheckAfterDefault default implementation of CheckAfter
func (input *Input) CheckAfterDefault() error {
	var err error
	input.Data.MainPort, input.Data.MainPortInt, err = Port(input.Data.MainPort)
	return err
}
