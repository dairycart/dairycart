# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L60))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L70))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L103))
- [x] Custom pagination (TestProductListRouteWithCustomFilter)

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L89))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L80))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L116))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L131))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L138))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L145))

Create product route:

- [x] Valid input ([TestProductCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L155))
- [ ] Invalid input
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L169))

Delete product route:

- [x] Newly Created SKU ([TestProductDeletionRouteForNewlyCreatedProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L181))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/33cde37591cf2cc441670bf099896348e050648f/integration_tests/main_test.go#L191))

## Product Attributes

Product attribute list route:

- [ ] Default pagination
- [ ] Custom pagination

Create product attribute route:

- [ ] Valid input
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
