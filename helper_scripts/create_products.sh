# new product
curl -X POST -d '{"name": "farts", "sku": "faaart", "product_weight": 1, "product_height": 2, "product_length": 3, "product_width": 4, "description": "this is a posted product", "upc": "12345"}' localhost:8080/product
# same product so we throw errors
curl -X POST -d '{"name": "farts", "sku": "faaart", "product_weight": 1, "product_height": 2, "product_length": 3, "product_width": 4, "description": "this is a posted product", "upc": "12345"}' localhost:8080/product
# another new product
curl -X POST -d '{"name": "farts", "sku": "fartz", "product_weight": 1, "product_height": 2, "product_length": 3, "product_width": 4, "description": "this is another posted product", "upc": "54321"}' localhost:8080/product