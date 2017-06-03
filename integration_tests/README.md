# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L75))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L85))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L117))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L128))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L105))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L95))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L143))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L156))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L167))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L175))

Create product route:

- [x] Valid input ([TestProductCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L186))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L210))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L198))

Delete product route:

- [x] Newly Created SKU ([TestProductDeletionRouteForNewlyCreatedProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L329))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L319))

## Product Attributes

Product attribute list route:

- [x] Default pagination ([TestProductAttributeListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L221))
- [x] Custom pagination ([TestProductAttributeListRetrievalWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L232))

Create product attribute route:

- [x] Valid input ([TestProductAttributeCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L247))
- [ ] Invalid input ([TestProductAttributeCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L259))
- [ ] Existent attribute name ([TestProductAttributeCreationWithAlreadyExistentName](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L270))

Update product attribute route:

- [x] Valid input ([TestProductAttributeUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L282))
- [x] Invalid input ([TestProductAttributeUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L294))
- [x] Nonexistent attribute ([TestProductAttributeUpdateForNonexistentAttribute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L305))

## Product Attribute Values

Create product attribute values route:

- [ ] Valid input
- [ ] Invalid input
- [ ] Existent attribute value

Update product attribute value route:

- [ ] Valid input
- [ ] Invalid input
- [ ] Nonexistent attribute
