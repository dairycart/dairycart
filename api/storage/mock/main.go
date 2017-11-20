package dairymock

import (
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}
