package models

import (
	"database/sql/driver"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	exampleSKU               = "example"
	exampleTimeString        = "2016-12-31 12:00:00.000000"
	exampleGarbageInput      = `{"things": "stuff"}`
	exampleMarshalTimeString = "2016-12-31T12:00:00.000000Z"
)

func TestDairytimeScan(t *testing.T) {
	t.Parallel()

	t.Run("with valid value", func(_t *testing.T) {
		_t.Parallel()

		test, err := time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
		require.Nil(t, err)

		dt := &Dairytime{}
		assert.NoError(t, dt.Scan(test))
	})

	t.Run("with invalid value", func(_t *testing.T) {
		_t.Parallel()

		dt := &Dairytime{}
		assert.Error(t, dt.Scan(nil))
	})
}

func TestDairytimeValue(t *testing.T) {
	t.Parallel()

	t.Run("with valid value", func(_t *testing.T) {
		_t.Parallel()

		test, err := time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
		require.Nil(t, err)
		dt := &Dairytime{test}

		expected := driver.Value(test)
		actual, err := dt.Value()
		assert.NoError(t, err)
		assert.Equal(t, expected, actual)
	})
}

func TestDairytimeMarshalText(t *testing.T) {
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

		out, err := time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
		require.Nil(t, err)

		expected := []byte(exampleMarshalTimeString)
		example := Dairytime{Time: out}
		actual, err := example.MarshalText()

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "Marshaled time string should marshal correctly")
	})

	t.Run("with zero time", func(_t *testing.T) {
		_t.Parallel()

		example := Dairytime{}
		actual, err := example.MarshalText()

		assert.Nil(t, err)
		assert.Nil(t, actual)
	})
}

func TestDairytimeUnmarshalText(t *testing.T) {
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

		example := []byte(exampleMarshalTimeString)
		nt := Dairytime{}
		err := nt.UnmarshalText(example)
		assert.Nil(t, err)
	})

	t.Run("with nil input", func(_t *testing.T) {
		_t.Parallel()

		nt := Dairytime{}
		err := nt.UnmarshalText(nil)
		assert.Nil(t, err)
	})

	t.Run("with empty byte slice", func(_t *testing.T) {
		_t.Parallel()

		nt := Dairytime{}
		err := nt.UnmarshalText([]byte{})
		assert.Nil(t, err)
	})
}

func TestDairytimeString(t *testing.T) {
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

		test, err := time.Parse("2006-01-02 03:04:00.000000", exampleTimeString)
		require.Nil(t, err)

		dt := &Dairytime{test}
		expected := "2016-12-31T12:00:00.000000Z"
		actual := dt.String()
		assert.Equal(t, expected, actual)
	})

	t.Run("with nil", func(_t *testing.T) {
		_t.Parallel()

		dt := (*Dairytime)(nil)
		expected := "nil"
		actual := dt.String()
		assert.Equal(t, expected, actual)
	})
}

func TestErrorResponseError(t *testing.T) {
	t.Parallel()

	t.Run("normal operation", func(_t *testing.T) {
		_t.Parallel()

		expected := "hi, mom"
		example := &ErrorResponse{Message: expected}

		actual := example.Error()
		assert.Equal(t, expected, actual)
	})
}
