package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiscountGenerateScanArgs(t *testing.T) {
	// this test will go away soon, but for now I can't stand having less than 100% test coverage
	t.Parallel()
	d := &Discount{}
	actual := d.generateScanArgs()
	assert.Equal(t, 8, len(actual), "there should be eight scan args for discounts")
}

func TestDiscountTypeIsValidWithValidInput(t *testing.T) {
	t.Parallel()
	d := &Discount{
		Type: "flat_discount",
	}
	assert.False(t, d.discountTypeIsValid())
}

func TestDiscountTypeIsValidWithInvalidInput(t *testing.T) {
	t.Parallel()
	d := &Discount{
		Type: "this is nonsense",
	}
	assert.False(t, d.discountTypeIsValid())
}
