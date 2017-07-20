# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L85))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L95))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L159))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L172))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L122))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L105))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L188))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L259))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L269))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L277))

Create product route:

- [x] Valid input ([TestProductCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L295))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L450))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L407))

Delete product route:

- [x] Newly created SKU ([TestProductDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L356))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L390))

## Product Options

Product option list route:

- [x] Default pagination ([TestProductOptionListRetrievalWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L460))
- [x] Custom pagination ([TestProductOptionListRetrievalWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L474))

Create product option route:

- [x] Valid input ([TestProductOptionCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L492))
- [x] Invalid input ([TestProductOptionCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L584))
- [x] Existent option name ([TestProductOptionCreationWithAlreadyExistentName](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L594))

Update product option route:

- [x] Valid input ([TestProductOptionUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L641))
- [x] Invalid input ([TestProductOptionUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L701))
- [x] Nonexistent option ([TestProductOptionUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L712))

Delete product option route:

- [] Newly created option value ([TestProductOptionDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L547))
- [] Nonexistent option ()

## Product Option Values

Create product option values route:

- [x] Valid input ([TestProductOptionValueCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L732))
- [x] Invalid input ([TestProductOptionValueCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L866))
- [x] Existent option value ([TestProductOptionValueCreationWithAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L876))

Update product option value route:

- [x] Valid input ([TestProductOptionValueUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L774))
- [x] Invalid input ([TestProductOptionValueUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L896))
- [x] Nonexistent option value ([TestProductOptionValueUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L907))
- [x] Duplicate option value ([TestProductOptionValueUpdateForAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L928))

Delete product option value route:

- [] Newly created option value ([TestProductOptionValueDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go.go#L830))
- [] Nonexistent option value ()

Single discount route:

- [x] Existent discount ([TestDiscountRetrievalForExistingDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L44))
- [x] Nonexistent discount ([TestDiscountRetrievalForNonexistentDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L67))

Discount list route:

- [x] Default pagination ([TestDiscountListRetrievalWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L86))
- [x] Custom pagination ([TestDiscountListRouteWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L99))

Discount creation route:

- [x] Valid input ([TestDiscountCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L115))
- [x] Invalid input ([TestDiscountCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L206))

Update discount route:

- [x] Valid input ([TestDiscountUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L217))
- [x] Invalid input ([TestDiscountUpdateInvalidDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L284))
- [x] Nonexistent Body ([TestDiscountUpdateWithInvalidBody](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L292))

Delete discount route:

- [] Newly created discount ([TestDiscountDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go.go#L168))
- [] Nonexistent discount ()

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
