INSERT INTO products
(
    "sku",
    "name",
    "description",
    "upc",
    "price",
    "cost",
    "quantity",
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
('t-shirt', 'T-Shirt', 'This is a t-shirt. wear it, or don''t. I''m not your dad.', null, 12.34, 5, 123, 1, 2, 3, 4, 5, 6, 7, 8),
('skateboard', 'Skateboard', 'This is a skateboard. Please wear a helmet.', '1234567890', 99.99, 50, 666, 8, 7, 6, 5, 4, 3, 2, 1),
('guitar', 'Guitar', 'It is a guitar', null, 6.66, 1.23, 321, 1, 2, 3, 4, 5, 6, 7, 8),
('fuzz-pedal', 'Fuzz Pedal', 'Make the guitar sound fuzzy', null, 256, 128, 100, 1, 2, 3, 4, 5, 6, 7, 8);

INSERT INTO product_options
(
    "name",
    "product_id"
)
VALUES
('color', 1),
('size', 1);

INSERT INTO product_option_values
(
    "product_option_id",
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
    "starts_on",
    "expires_on"
)
VALUES
(
    '10% off',
    'percentage',
    10.00,
    NOW(),
    NOW() + (1 * interval '1 month')
),
(
    '50% off',
    'percentage',
    50.00,
    NOW(),
    null
),
(
    'New customer special',
    'flat_amount',
    10.00,
    NOW(),
    null
);