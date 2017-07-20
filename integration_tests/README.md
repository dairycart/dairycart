# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L85-L93))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L95-L103))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L159-L170))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L172-L186))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L122-L157))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L105-L120))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L188-L257))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L259-L267))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L269-L275))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L277-L293))

Create product route:

- [x] Valid input ([TestProductCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L295-L354))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L450-L458))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L407-L448))

Delete product route:

- [x] Newly created SKU ([TestProductDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L356-L388))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L390-L405))

## Product Options

Product option list route:

- [x] Default pagination ([TestProductOptionListRetrievalWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L460-L472))
- [x] Custom pagination ([TestProductOptionListRetrievalWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L474-L490))

Create product option route:

- [x] Valid input ([TestProductOptionCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L492-L545))
- [x] Invalid input ([TestProductOptionCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L584-L592))
- [x] Existent option name ([TestProductOptionCreationWithAlreadyExistentName](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L594-L639))

Update product option route:

- [x] Valid input ([TestProductOptionUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L641-L699))
- [x] Invalid input ([TestProductOptionUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L701-L710))
- [x] Nonexistent option ([TestProductOptionUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L712-L730))

Delete product option route:

- [] Newly created option value ([TestProductOptionDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L547-L582))
- [] Nonexistent option ()

## Product Option Values

Create product option values route:

- [x] Valid input ([TestProductOptionValueCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L732-L772))
- [x] Invalid input ([TestProductOptionValueCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L866-L874))
- [x] Existent option value ([TestProductOptionValueCreationWithAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L876-L894))

Update product option value route:

- [x] Valid input ([TestProductOptionValueUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L774-L828))
- [x] Invalid input ([TestProductOptionValueUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L896-L905))
- [x] Nonexistent option value ([TestProductOptionValueUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L907-L926))
- [x] Duplicate option value ([TestProductOptionValueUpdateForAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L928-L945))

Delete product option value route:

- [] Newly created option value ([TestProductOptionValueDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L830-L864))
- [] Nonexistent option value ()

Single discount route:

- [x] Existent discount ([TestDiscountRetrievalForExistingDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L44-L65))
- [x] Nonexistent discount ([TestDiscountRetrievalForNonexistentDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L67-L84))

Discount list route:

- [x] Default pagination ([TestDiscountListRetrievalWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L86-L97))
- [x] Custom pagination ([TestDiscountListRouteWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L99-L113))

Discount creation route:

- [x] Valid input ([TestDiscountCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L115-L166))
- [x] Invalid input ([TestDiscountCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L206-L215))

Update discount route:

- [x] Valid input ([TestDiscountUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L217-L282))
- [x] Invalid input ([TestDiscountUpdateInvalidDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L284-L290))
- [x] Nonexistent Body ([TestDiscountUpdateWithInvalidBody](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L292-L297))

Delete discount route:

- [] Newly created discount ([TestDiscountDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L168-L204))
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
