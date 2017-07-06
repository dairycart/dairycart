# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L17))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L27))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L61))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L72))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L48))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L37))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L87))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L100))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L110))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L118))

Create product route:

- [x] Valid input ([TestProductCreationAndDeletion](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L349))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L142))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L130))

Delete product route:

- [x] Newly created SKU ([TestProductCreationAndDeletion](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L349))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L338))

## Product Options

Product option list route:

- [x] Default pagination ([TestProductOptionListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L152))
- [x] Custom pagination ([TestProductOptionListRetrievalWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L163))

Create product option route:

- [x] Valid input ([TestProductOptionCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L178))
- [x] Invalid input ([TestProductOptionCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L190))
- [x] Existent option name ([TestProductOptionCreationWithAlreadyExistentName](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L200))

Update product option route:

- [x] Valid input ([TestProductOptionUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L214))
- [x] Invalid input ([TestProductOptionUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L226))
- [x] Nonexistent option ([TestProductOptionUpdateForNonexistentOption](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L237))

Delete product option route:

- [] Newly created option value ()
- [] Nonexistent option ()

## Product Option Values

Create product option values route:

- [x] Valid input ([TestProductOptionValueCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L250))
- [x] Invalid input ([TestProductOptionValueCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L262))
- [x] Existent option value ([TestProductOptionValueCreationWithAlreadyExistentValue](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L272))

Update product option value route:

- [x] Valid input ([TestProductOptionValueUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L285))
- [x] Invalid input ([TestProductOptionValueUpdateWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L297))
- [x] Nonexistent option value ([TestProductOptionValueUpdateForNonexistentOption](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L308))
- [x] Duplicate option value ([TestProductOptionValueUpdateForAlreadyExistentValue](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/products_test.go.go#L321))

Delete product option value route:

- [] Newly created option value ()
- [] Nonexistent option value ()

Single discount route:

- [x] Existent discount ([TestDiscountRetrievalForExistingDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L17))
- [x] Nonexistent discount ([TestDiscountRetrievalForNonexistentDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L31))

Discount list route:

- [x] Default pagination ([TestDiscountListRetrievalWithDefaultFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L44))
- [x] Custom pagination ([TestDiscountListRouteWithCustomFilter](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L56))

Discount creation route:

- [x] Valid input ([TestDiscountCreation](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L72))
- [x] Invalid input ([TestDiscountCreationWithInvalidInput](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L84))

Update discount route:

- [x] Valid input ([TestDiscountUpdate](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L95))
- [x] Invalid input ([TestDiscountUpdateInvalidDiscount](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L107))
- [x] Nonexistent Body ([TestDiscountUpdateWithInvalidBody](https://github.com/verygoodsoftwarenotvirus/dairycart/blob/master/integration_tests/pricing_test.go.go#L115))

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

User login scenarios:

- [] Create user and log out ()
- [] Create user, log out, log back in, and log out again ()
