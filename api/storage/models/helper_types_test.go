package models

import (
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

const (
	exampleSKU = "example"
	//exampleTimeString        = "2016-12-01 12:00:00.000000"
	exampleGarbageInput      = `{"things": "stuff"}`
	exampleMarshalTimeString = "2016-12-31T12:00:00.000000Z"
)

func TestNullTimeMarshalText(t *testing.T) {
	t.Parallel()

	out, err := time.Parse("2006-01-02 03:04:00.000000", "2016-12-31 12:00:00.000000")
	require.Nil(t, err)

	expected := []byte(exampleMarshalTimeString)
	example := NullTime{pq.NullTime{Time: out, Valid: true}}
	actual, err := example.MarshalText()

	require.Nil(t, err)
	require.Equal(t, expected, actual, "Marshaled time string should marshal correctly")
}

func TestNullTimeUnmarshalText(t *testing.T) {
	t.Parallel()
	example := []byte(exampleMarshalTimeString)
	nt := NullTime{}
	err := nt.UnmarshalText(example)
	require.Nil(t, err)
}
