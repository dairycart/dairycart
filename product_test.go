package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProductExistsInDB(t *testing.T) {
	productExists, err := ProductExistsInDB("example_sku")
	if err != nil {
		t.FailNow()
	}

	assert.True(t, productExists, "product that should exist is being returned as non-existent")
}
