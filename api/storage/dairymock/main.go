package dairymock

import (
// "github.com/stretchr/testify/mock"
)

type ReturnInstruction struct {
	// mock.Mock
	// CallLimit is the number of times a function should be called before returning an error
	CallLimit int
	// CallCount is the number of times a given function has been called
	CallCount int
	// ShouldReturnError indicates whether an error should be returned after a given number of calls
	ShouldReturnError bool
	// Error is the error we should return in the event we're instructed to do so
	Error error
}

type MockDB struct {
	InstructionMap map[string]ReturnInstruction
}
