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

- [x] Newly Created SKU ([TestProductDeletionRouteForNewlyCreatedProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L456))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L445))

## Product Attributes

Product attribute list route:

- [x] Default pagination ([TestProductAttributeListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L222))
- [x] Custom pagination ([TestProductAttributeListRetrievalWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L233))

Create product attribute route:

- [x] Valid input ([TestProductAttributeCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L248))
- [ ] Invalid input ([TestProductAttributeCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L260))
- [ ] Existent attribute name ([TestProductAttributeCreationWithAlreadyExistentName](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L271))

Update product attribute route:

- [x] Valid input ([TestProductAttributeUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L283))
- [x] Invalid input ([TestProductAttributeUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L295))
- [x] Nonexistent attribute ([TestProductAttributeUpdateForNonexistentAttribute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L306))

## Product Attribute Values

Create product attribute values route:

- [x] Valid input ([TestProductAttributeValueCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L319))
- [x] Invalid input ([TestProductAttributeValueCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L331))
- [x] Existent attribute value ([TestProductAttributeValueCreationWithAlreadyExistentValue](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L342))

Update product attribute value route:

- [x] Valid input ([TestProductAttributeValueUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L354))
- [x] Invalid input ([TestProductAttributeValueUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L366))
- [x] Nonexistent attribute value ([TestProductAttributeValueUpdateForNonexistentAttribute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L378))
- [x] Duplicate attribute value ([TestProductAttributeValueUpdateForAlreadyExistentValue](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L391))

Single discount route:

- [x] Existent discount ([TestDiscountRetrievalForExistingDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L419))
- [x] Nonexistent discount ([TestDiscountRetrievalForNonexistentDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L431))
