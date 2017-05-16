package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductExistsInDB(t *testing.T) {
	fakeDB := NewMockDB()

	_, err := productExistsInDB(fakeDB, "example_sku")
	if err != nil {
		t.FailNow()
	}

	assert.Equal(t, fakeDB.CallList, []string{"Model", "QueryOne"})
}
