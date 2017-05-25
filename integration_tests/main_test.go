package dairytest

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	err := ensureThatDairycartIsAlive()
	if err != nil {
		log.Fatalf("dairycart isn't up: %v", err)
	}
}

func TestTesting(t *testing.T) {
	assert.True(t, true)
}
