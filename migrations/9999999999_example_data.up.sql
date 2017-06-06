INSERT INTO product_progenitors
(
    "name",
    "description",
    "price",
    "cost",
    "product_weight",
    "product_height",
    "product_width",
    "product_length",
    "package_weight",
    "package_height",
    "package_width",
    "package_length"
)
VALUES
('T-Shirt', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', 12.34, 5, 1, 2, 3, 4, 5, 6, 7, 8),
('Skateboard', 'This is a skateboard. Please wear a helmet.', 99.99, 50, 8, 7, 6, 5, 4, 3, 2, 1);

INSERT INTO products
(
    "product_progenitor_id",
    "sku",
    "name",
    "upc",
    "price",
    "cost",
    "quantity"
)
VALUES
(1, 't-shirt-small-red', 'Red T-Shirt (Small)', null, 12.34, 5, 123),
(1, 't-shirt-small-green', 'Green T-Shirt (Small)', null, 12.34, 5, 123),
(1, 't-shirt-small-blue', 'Blue T-Shirt (Small)', null, 12.34, 5, 123),
(1, 't-shirt-medium-red', 'Red T-Shirt (Medium)', null, 12.34, 5, 123),
(1, 't-shirt-medium-green', 'Green T-Shirt (Medium)', null, 12.34, 5, 123),
(1, 't-shirt-medium-blue', 'Blue T-Shirt (Medium)', null, 12.34, 5, 123),
(1, 't-shirt-large-red', 'Red T-Shirt (Large)', null, 12.34, 5, 123),
(1, 't-shirt-large-green', 'Green T-Shirt (Large)', null, 12.34, 5, 123),
(1, 't-shirt-large-blue', 'Blue T-Shirt (Large)', null, 12.34, 5, 123),
(2, 'skateboard', 'Skateboard', '1234567890', 99.99, 50, 123);

INSERT INTO product_attributes
(
    "name",
    "product_progenitor_id"
)
VALUES
('color', 1),
('size', 1);

INSERT INTO product_attribute_values
(
    "product_attribute_id",
    "value"
)
VALUES
(1, 'red'),
(1, 'green'),
(1, 'blue'),
(2, 'small'),
(2, 'medium'),
(2, 'large');

INSERT INTO discounts
(
    "name",
    "type",
    "amount",
    "product_id",
    "starts_on",
    "expires_on"
)
VALUES
(
    '10% off',
    'percentage',
    10.00,
    10,
    NOW(),
    NOW() + (1 * interval '1 month')
);