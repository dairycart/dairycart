# Integration tests checklist

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L85-L93))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L95-L103))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L158-L169))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L171-L185))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L121-L156))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L105-L119))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L187-L256))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L258-L266))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L268-L274))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L276-L291))

Create product route:

- [x] Valid input ([TestProductCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L293-L352))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L446-L454))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L404-L444))

Delete product route:

- [x] Newly created SKU ([TestProductDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L354-L386))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L388-L402))

## Product Options

Product option list route:

- [x] Default pagination ([TestProductOptionListRetrievalWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L456-L468))
- [x] Custom pagination ([TestProductOptionListRetrievalWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L470-L486))

Create product option route:

- [x] Valid input ([TestProductOptionCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L488-L541))
- [x] Invalid input ([TestProductOptionCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L592-L600))
- [x] Existent option name ([TestProductOptionCreationWithAlreadyExistentName](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L602-L647))

Update product option route:

- [x] Valid input ([TestProductOptionUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L649-L707))
- [x] Invalid input ([TestProductOptionUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L709-L718))
- [x] Nonexistent option ([TestProductOptionUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L720-L737))

Delete product option route:

- [x] Newly created option value ([TestProductOptionDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L543-L578))
- [x] Nonexistent option ([TestProductOptionDeletionForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L580-L590))

## Product Option Values

Create product option values route:

- [x] Valid input ([TestProductOptionValueCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L739-L779))
- [x] Invalid input ([TestProductOptionValueCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L885-L893))
- [x] Existent option value ([TestProductOptionValueCreationWithAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L895-L912))

Update product option value route:

- [x] Valid input ([TestProductOptionValueUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L781-L835))
- [x] Invalid input ([TestProductOptionValueUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L914-L923))
- [x] Nonexistent option value ([TestProductOptionValueUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L925-L943))
- [x] Duplicate option value ([TestProductOptionValueUpdateForAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L945-L962))

Delete product option value route:

- [x] Newly created option value ([TestProductOptionValueDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L837-L871))
- [x] Nonexistent option value ([TestProductOptionValueDeletionForNonexistentOptionValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L873-L883))

Single discount route:

- [x] Existent discount ([TestDiscountRetrievalForExistingDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L44-L65))
- [x] Nonexistent discount ([TestDiscountRetrievalForNonexistentDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L67-L83))

Discount list route:

- [x] Default pagination ([TestDiscountListRetrievalWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L85-L96))
- [x] Custom pagination ([TestDiscountListRouteWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L98-L112))

Discount creation route:

- [x] Valid input ([TestDiscountCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L114-L163))
- [x] Invalid input ([TestDiscountCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L215-L224))

Update discount route:

- [x] Valid input ([TestDiscountUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L226-L289))
- [x] Invalid input ([TestDiscountUpdateInvalidDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L291-L297))
- [x] Nonexistent Body ([TestDiscountUpdateWithInvalidBody](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L299-L304))

Delete discount route:

- [x] Newly created discount ([TestDiscountDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L165-L201))
- [x] Nonexistent discount ([TestDiscountDeletionForNonexistentDiscount](https://github.com/dairycart/dairycart/blob/master/integration_tests/pricing_test.go#L203-L213))

User creation routes:

- [x] Valid input, regular user ([TestUserCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L26-L58))
- [x] Valid input, admin user ([TestAdminUserCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L88-L120))
- [x] Valid input, admin user, without admin credentials ()
- [x] Already existent user ([TestUserCreationForAlreadyExistentUsername](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L122-L165))
- [x] Invalid password ([TestUserCreationWithInvalidPassword](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L60-L74))
- [x] Invalid input ([TestUserCreationWithInvalidCreationBody](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L76-L86))

User deletion routes:

- [] Regular user ()
- [] Admin user ()
- [] Admin user without admin credentials ()
- [] Non existent user ()

User login routes:

- [] Valid password ()
- [] Non-existent user ()
- [] Invalid input ()
- [] Bad password ()

User logout routes:

- [] Log out ()
