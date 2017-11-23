// +build !migrated

package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRoute(t *testing.T) {
	inputOutputMap := []struct {
		In  []string
		Out string
	}{
		{
			In:  []string{"tests"},
			Out: "/v1/tests",
		},
		{
			In:  []string{"things", "and", "stuff"},
			Out: "/v1/things/and/stuff",
		},
	}

	for _, x := range inputOutputMap {
		input := x.In
		expected := x.Out
		actual := buildRoute("v1", input...)
		assert.Equal(t, expected, actual, `buildRoute with input of ["%s"] should equal %s`, strings.Join(input, `", "`), expected)
	}
}
