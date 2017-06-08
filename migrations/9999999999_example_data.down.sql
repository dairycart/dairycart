DELETE FROM variants WHERE id IS NOT NULL;
DELETE FROM products WHERE id IS NOT NULL;
DELETE FROM product_options WHERE id IS NOT NULL;
DELETE FROM product_option_values WHERE id IS NOT NULL;
DELETE FROM product_attributes WHERE id IS NOT NULL;
DELETE FROM product_attribute_values WHERE id IS NOT NULL;
DELETE FROM discounts WHERE id IS NOT NULL;