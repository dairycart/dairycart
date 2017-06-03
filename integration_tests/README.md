# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L72))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L82))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L114))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L125))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L102))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L92))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L140))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L154))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L162))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L170))

Create product route:

- [x] Valid input ([TestProductCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L181))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L205))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L193))

Delete product route:

- [x] Newly Created SKU ([TestProductDeletionRouteForNewlyCreatedProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L266))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L256))

## Product Attributes

Product attribute list route:

- [x] Default pagination ([TestProductAttributeListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L216))
- [x] Custom pagination ([TestProductAttributeListRetrievalWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L227))

Create product attribute route:

- [x] Valid input ([TestProductAttributeCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/main_test.go#L242))
- [ ] Invalid input
- [ ] Existent attribute name

Update product attribute route:

- [ ] Valid input
- [ ] Invalid input
- [ ] Nonexistent attribute

## Product Attribute Values

Create product attribute values route:

- [ ] Valid input
- [ ] Invalid input
- [ ] Existent attribute value

Update product attribute value route:

- [ ] Valid input
- [ ] Invalid input
- [ ] Nonexistent attribute
