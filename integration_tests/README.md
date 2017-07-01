# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L10))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L20))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L53))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L64))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L41))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L30))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L79))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L92))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L103))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L111))

Create product route:

- [x] Valid input ([TestProductCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L123))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L147))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L135))

Delete product route:

- [x] Newly Created SKU ([TestProductDeletionRouteForNewlyCreatedProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L357))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L346))

## Product Options

Product option list route:

- [x] Default pagination ([TestProductOptionListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L158))
- [x] Custom pagination ([TestProductOptionListRetrievalWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L169))

Create product option route:

- [x] Valid input ([TestProductOptionCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L184))
- [x] Invalid input ([TestProductOptionCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L196))
- [x] Existent option name ([TestProductOptionCreationWithAlreadyExistentName](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L207))

Update product option route:

- [x] Valid input ([TestProductOptionUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L219))
- [x] Invalid input ([TestProductOptionUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L231))
- [x] Nonexistent option ([TestProductOptionUpdateForNonexistentOption](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L242))

## Product Option Values

Create product option values route:

- [x] Valid input ([TestProductOptionValueCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L255))
- [x] Invalid input ([TestProductOptionValueCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L267))
- [x] Existent option value ([TestProductOptionValueCreationWithAlreadyExistentValue](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L278))

Update product option value route:

- [x] Valid input ([TestProductOptionValueUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L290))
- [x] Invalid input ([TestProductOptionValueUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L302))
- [x] Nonexistent option value ([TestProductOptionValueUpdateForNonexistentOption](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L314))
- [x] Duplicate option value ([TestProductOptionValueUpdateForAlreadyExistentValue](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L327))

Single discount route:

- [x] Existent discount ([TestDiscountRetrievalForExistingDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L16))
- [x] Nonexistent discount ([TestDiscountRetrievalForNonexistentDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L28))

Discount list route:

- [x] Default pagination ([TestDiscountListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L40))
- [x] Custom pagination ([TestDiscountListRouteWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L52))

Discount creation route:

- [x] Valid input ([TestDiscountCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L68))
- [x] Invalid input ([TestDiscountCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L80))

Update discount route:

- [x] Valid input ([TestDiscountUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L91))
- [x] Invalid input ([TestDiscountUpdateInvalidDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L103))
- [x] Nonexistent Body ([TestDiscountUpdateWithInvalidBody](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L110))

User creation routes:

- [] Valid input ()
- [] Already existent user ()
- [] Invalid input ()
- [] Bad password ()

User login routes:

- [] Valid password ()
- [] Non-existent user ()
- [] Invalid input ()
- [] Bad password ()