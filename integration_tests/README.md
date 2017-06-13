# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L74))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L84))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L117))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L128))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L105))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L94))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L143))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L156))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L167))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L175))

Create product route:

- [x] Valid input ([TestProductCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L187))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L211))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L199))

Delete product route:

- [x] Newly Created SKU ([TestProductDeletionRouteForNewlyCreatedProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L507))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L496))

## Product Options

Product option list route:

- [x] Default pagination ([TestProductOptionListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L222))
- [x] Custom pagination ([TestProductOptionListRetrievalWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L233))

Create product option route:

- [x] Valid input ([TestProductOptionCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L248))
- [x] Invalid input ([TestProductOptionCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L260))
- [x] Existent option name ([TestProductOptionCreationWithAlreadyExistentName](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L271))

Update product option route:

- [x] Valid input ([TestProductOptionUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L283))
- [x] Invalid input ([TestProductOptionUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L295))
- [x] Nonexistent option ([TestProductOptionUpdateForNonexistentOption](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L306))

## Product Option Values

Create product option values route:

- [x] Valid input ([TestProductOptionValueCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L319))
- [x] Invalid input ([TestProductOptionValueCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L331))
- [x] Existent option value ([TestProductOptionValueCreationWithAlreadyExistentValue](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L342))

Update product option value route:

- [x] Valid input ([TestProductOptionValueUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L354))
- [x] Invalid input ([TestProductOptionValueUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L366))
- [x] Nonexistent option value ([TestProductOptionValueUpdateForNonexistentOption](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L378))
- [x] Duplicate option value ([TestProductOptionValueUpdateForAlreadyExistentValue](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L391))

Single discount route:

- [x] Existent discount ([TestDiscountRetrievalForExistingDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L419))
- [x] Nonexistent discount ([TestDiscountRetrievalForNonexistentDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L431))

Discount list route:

- [x] Default pagination ([TestDiscountListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L443))
- [x] Custom pagination ([TestDiscountListRouteWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L455))

Discount creation route:

- [x] Valid input ([TestDiscountCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L471))
- [x] Invalid input ([TestDiscountCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L483))
