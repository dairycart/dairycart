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
- [x] Invalid input ([TestProductOptionCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L596-L604))
- [x] Existent option name ([TestProductOptionCreationWithAlreadyExistentName](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L606-L651))

Update product option route:

- [x] Valid input ([TestProductOptionUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L653-L711))
- [x] Invalid input ([TestProductOptionUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L713-L722))
- [x] Nonexistent option ([TestProductOptionUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L724-L742))

Delete product option route:

- [x] Newly created option value ([TestProductOptionDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L547-L582))
- [x] Nonexistent option ([TestProductOptionDeletionForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L584-L594))

## Product Option Values

Create product option values route:

- [x] Valid input ([TestProductOptionValueCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L744-L784))
- [x] Invalid input ([TestProductOptionValueCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L890-L898))
- [x] Existent option value ([TestProductOptionValueCreationWithAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L900-L918))

Update product option value route:

- [x] Valid input ([TestProductOptionValueUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L786-L840))
- [x] Invalid input ([TestProductOptionValueUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L920-L929))
- [x] Nonexistent option value ([TestProductOptionValueUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L931-L950))
- [x] Duplicate option value ([TestProductOptionValueUpdateForAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L952-L969))

Delete product option value route:

- [x] Newly created option value ([TestProductOptionValueDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L842-L876))
- [x] Nonexistent option value ([TestProductOptionValueDeletionForNonexistentOptionValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L878-L888))

Single discount route:

- [x] Existent discount ([TestDiscountRetrievalForExistingDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L44-L65))
- [x] Nonexistent discount ([TestDiscountRetrievalForNonexistentDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L67-L84))

Discount list route:

- [x] Default pagination ([TestDiscountListRetrievalWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L86-L97))
- [x] Custom pagination ([TestDiscountListRouteWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L99-L113))

Discount creation route:

- [x] Valid input ([TestDiscountCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L115-L166))
- [x] Invalid input ([TestDiscountCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L218-L227))

Update discount route:

- [x] Valid input ([TestDiscountUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L229-L294))
- [x] Invalid input ([TestDiscountUpdateInvalidDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L296-L302))
- [x] Nonexistent Body ([TestDiscountUpdateWithInvalidBody](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L304-L309))

Delete discount route:

- [x] Newly created discount ([TestDiscountDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L168-L204))
- [x] Nonexistent discount ([TestDiscountDeletionForNonexistentDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L206-L216))

User creation routes:

- [x] Valid input, regular user (TestUserCreation)
- [x] Valid input, admin user (TestAdminUserCreation)
- [x] Already existent user (TestUserCreationForAlreadyExistentUsername)
- [x] Invalid password (TestUserCreationWithInvalidPassword)
- [x] Invalid input (TestUserCreationWithInvalidCreationBody)

User deletion routes:

- [] Regular user ()
- [] Admin user ()
- [] Non existent user ()

User login routes:

- [] Valid password ()
- [] Non-existent user ()
- [] Invalid input ()
- [] Bad password ()

User login scenarios:

- [] Create user and log out ()
- [] Create user, log out, log back in, and log out again ()
