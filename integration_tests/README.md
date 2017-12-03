# Integration tests checklist

## WARNING

The plan is for this document to be valid very soon, but switching to a better testing method/structure means that this document will not be updated, nor will its quality be ensured in tests until this notice disappears.

## Product routes

Product existence route:

- [x] Existent SKU ([TestProductExistenceRouteForExistingProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L87-L95))
- [x] Nonexistent SKU ([TestProductExistenceRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L97-L105))

Product list route:

- [x] Default pagination ([TestProductListRouteWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L161-L172))
- [x] Custom pagination ([TestProductListRouteWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L174-L188))

Single product route:

- [x] Existent SKU ([TestProductRetrievalRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L123-L159))
- [x] Nonexistent SKU ([TestProductRetrievalRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L107-L121))

Update product route:

- [x] Valid input ([TestProductUpdateRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L190-L263))
- [x] Invalid  ([TestProductUpdateRouteWithCompletelyInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L265-L273))
- [x] Invalid SKU input ([TestProductUpdateRouteWithInvalidSKU](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L275-L281))
- [x] Nonexistent SKU ([TestProductUpdateRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L283-L298))

Create product route:

- [x] Valid input ([TestProductCreationRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L300-L382))
- [x] Invalid input ([TestProductCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1226-L1234))
- [x] Existent SKU ([TestProductCreationWithAlreadyExistentSKU](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1184-L1224))
- [x] Product with multiple options ([TestProductCreationRouteWithOptions](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L384-L775))

Delete product route:

- [x] Newly created SKU ([TestProductDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L777-L808))
- [x] Nonexistent SKU ([TestProductDeletionRouteForNonexistentProduct](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L810-L824))

## Product Root Routes

Product Root list route:

- [x] Default pagination ([TestProductRootListRetrievalRouteWithDefaultPagination](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L826-L837))
- [x] Custom pagination ([TestProductRootListRetrievalRouteWithCustomPagination](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L839-L854))

Single product root route:

- [x] Existent root ([TestProductRootRetrievalRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L856-L1122))
- [x] Nonexistent root ([TestProductRootRetrievalRouteForNonexistentRoot](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1124-L1133))

Delete product root route:

- [x] Existent root ([TestProductRootDeletionRoute](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1135-L1171))
- [x] Nonexistent root ([TestProductRootDeletionRouteForNonexistentRoot](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1173-L1182))

## Product Options

Product option list route:

- [x] Default pagination ([TestProductOptionListRetrievalWithDefaultFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1236-L1248))
- [x] Custom pagination ([TestProductOptionListRetrievalWithCustomFilter](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1250-L1266))

Create product option route:

- [x] Valid input ([TestProductOptionCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1268-L1319))
- [x] Invalid input ([TestProductOptionCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1370-L1378))
- [x] Existent option name ([TestProductOptionCreationWithAlreadyExistentName](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1380-L1443))

Update product option route:

- [x] Valid input ([TestProductOptionUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1445-L1534))
- [x] Invalid input ([TestProductOptionUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1536-L1545))
- [x] Nonexistent option ([TestProductOptionUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1547-L1564))

Delete product option route:

- [x] Newly created option value ([TestProductOptionDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1321-L1356))
- [x] Nonexistent option ([TestProductOptionDeletionForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1358-L1368))

## Product Option Values

Create product option values route:

- [x] Valid input ([TestProductOptionValueCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1566-L1605))
- [x] Invalid input ([TestProductOptionValueCreationWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1763-L1771))
- [x] Existent option value ([TestProductOptionValueCreationWithAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1773-L1790))

Update product option value route:

- [x] Valid input ([TestProductOptionValueUpdate](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1607-L1713))
- [x] Invalid input ([TestProductOptionValueUpdateWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1792-L1801))
- [x] Nonexistent option value ([TestProductOptionValueUpdateForNonexistentOption](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1803-L1821))
- [x] Duplicate option value ([TestProductOptionValueUpdateForAlreadyExistentValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1823-L1837))

Delete product option value route:

- [x] Newly created option value ([TestProductOptionValueDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1715-L1749))
- [x] Nonexistent option value ([TestProductOptionValueDeletionForNonexistentOptionValue](https://github.com/dairycart/dairycart/blob/master/integration_tests/products_test.go#L1751-L1761))

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

- [x] Valid input, regular user ([TestUserCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L27-L60))
- [x] Valid input, admin user ([TestAdminUserCreation](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L92-L125))
- [x] Valid input, admin user, without admin credentials ([TestAdminUserCreationFailsWithoutAdminCredentials](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L127-L134))
- [x] Already existent user ([TestUserCreationForAlreadyExistentUsername](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L136-L180))
- [x] Invalid password ([TestUserCreationWithInvalidPassword](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L62-L77))
- [x] Invalid input ([TestUserCreationWithInvalidCreationBody](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L79-L90))

User deletion routes:

- [x] Regular user ([TestUserDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L182-L214))
- [x] Regular user without admin credentials ([TestUserDeletionAsRegularUser](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L228-L262))
- [x] Admin user ([TestAdminUserDeletion](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L264-L296))
- [x] Admin user without admin credentials ([TestAdminUserDeletionAsRegularUser](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L298-L332))
- [x] Non existent user ([TestUserDeletionForNonexistentUser](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L216-L226))

User login routes:

- [x] Valid password ([TestUserLogin](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L334-L391))
- [x] Bad password ([TestUserLoginWithInvalidPassword](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L393-L450))
- [x] Nonexistent user ([TestUserLoginForNonexistentUser](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L463-L471))
- [x] Invalid input ([TestUserLoginWithInvalidInput](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L452-L461))

User logout routes:

- [x] Valid user ([TestUserLogout](https://github.com/dairycart/dairycart/blob/master/integration_tests/auth_test.go#L473-L521))
