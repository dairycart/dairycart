package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	skuValidationPattern = `^[a-zA-Z\-_]+$`
)

var skuValidator *regexp.Regexp

func init() {
	skuValidator = regexp.MustCompile(skuValidationPattern)
}

//////////////////////////////////////////////////////////////////////////////////
//                                             _\/                              //
//                                           .-'.'`)                            //
//    Products!                           .-' .'                                //
//                  .                  .-'     `-.          __\/                //
//                   \.    .  |,   _.-'       -:````)   _.-'.'``)               //
//                    \`.  |\ | \.-_.           `._ _.-' .'`                    //
//                   __) )__\ |! )/ \_.          _.-'      `.                   //
//               _.-'__`-' =`:' /.' / |      _.-'      -:`````)                 //
//         __.--' ( (@> ))  = \ ^ `'. |_. .-'           `.                      //
//        : @       `^^^    == \ ^   `. |<                `.                    //
//        VvvvvvvvVvvvv)    =  ;   ^  ;_/ :           -:``````)                 //
//          (^^^^^^^^^^=  ==   |      ; \. :            `.                      //
//       ((  `-----------.  == |  ^   ;_/   :             `.                    //
//       /\             /==   /       : \.  :     _..--``````)                  //
//     __\ \_          ; ==  /  ^     :_/   :      `.                           //
//   ><__   _```---.._/ ====/       ^ : \. :         `.                         //
//       / / `._  ^  ;==   /  ^        :/ .            `.                       //
//       \/ ((  `._ / === /       ^    `.'       _.--`````)                     //
//       (( /\     ;===  /      ^     .'        `.                              //
//        __\ \_  : === | ^ /                     `.                            //
//     >><__   _``--...__.-'   ^  /  ^              `.                          //
//          / / `._        ^    .'              .--`````)     .--..             //
//          \/   :=`--...____.-'  ^     .___.-'|            .' .--.`.   (       //
//         ((    | ===    \                  `.|__.         ; ^:   `.'   )      //
//                :  ====  \  ^      ^         `. |          ;  `;   `../__     //
//              .-'\====    \    .       ^      `.|___.       ;^  `;       \    //
//           .-'    :  ===   \.-'              ^  `.  |        ;  ^ `;      )   //
//        .-'    ^   \==== .-'         ^            `.|___.     ;     ;    (    //
//     .-'         ^  :=.-'                 ^        `.   |      ;     ;        //
//   .'      ^       .-'          ^               ^    ;_/__.    ;  ^   ;       //
//   : ^        ^ .-'     ^                   ^   ;     ;   |    ;       ;      //
//   :    ^    .-'    ^          ^      ^     _.-'    ^  ;_/    ; ^      ;      //
//   :  ^    .'                           _.-"    ^      ; \.  ;      ^  ;      //
//    `.   ^ :   ^         ^       ^__.--"               ;_/  ;          ;      //
//      `.^  :                __.--"\         ^     ^    ; \ ;     ^     ;      //
//        `-.:       ^___.---"\ ===  \   ^               ;_/'        ^  ;       //
//           ``.^         `.   `\===  \         ^     ^       ^        ;        //
//              `.     ^    `.   `-. ==\                          ^   ;         //
//               _`-._        `.    `\= \    ^           ^           ;          //
//       __..--''    _`-._^     `.    `-.`\         ^          ^    ;           //
//      (-(-(-(-(--''     `-._  ^ `.     `\`\              ^      .'            //
//                    __..---''     `._     `-. ^      ^      ^ .'              //
//           __..---''    ___....---'`-`)      `---...____..---'                //
//          (-(-(-(-(---''             '                                        //
//////////////////////////////////////////////////////////////////////////////////

// Product describes something a user can buy
type Product struct {
	// Basic Info
	ID                  int64      `json:"id"`
	ProductProgenitorID int64      `json:"product_progenitor_id"`
	SKU                 string     `json:"sku"`
	Name                string     `json:"name"`
	UPC                 NullString `json:"upc"`
	Quantity            int        `json:"quantity"`

	// Pricing Fields
	OnSale    bool        `json:"on_sale"`
	Price     float32     `json:"price"`
	SalePrice NullFloat64 `json:"sale_price"`

	// // Housekeeping
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  pq.NullTime `json:"-"`
	ArchivedAt pq.NullTime `json:"-"`

	ProductProgenitor
}

// generateScanArgs generates an array of pointers to struct fields for sql.Scan to populate
func (p *Product) generateScanArgs() []interface{} {
	return []interface{}{
		&p.ID,
		&p.ProductProgenitorID,
		&p.SKU,
		&p.Name,
		&p.UPC,
		&p.Quantity,
		&p.OnSale,
		&p.Price,
		&p.SalePrice,
		&p.CreatedAt,
		&p.UpdatedAt,
		&p.ArchivedAt,
	}
}

// generateJoinScanArgs does some stuff TODO: write better docs
func (p *Product) generateJoinScanArgs() []interface{} {
	productScanArgs := p.generateScanArgs()
	progenitorScanArgs := p.ProductProgenitor.generateScanArgs()
	return append(productScanArgs, progenitorScanArgs...)
}

// ProductsResponse is a product response struct
type ProductsResponse struct {
	ListResponse
	Data []Product `json:"data"`
}

// ProductCreationInput is a struct that represents a product creation body
type ProductCreationInput struct {
	Description         string                           `json:"description"`
	Taxable             bool                             `json:"taxable"`
	ProductWeight       float32                          `json:"product_weight"`
	ProductHeight       float32                          `json:"product_height"`
	ProductWidth        float32                          `json:"product_width"`
	ProductLength       float32                          `json:"product_length"`
	PackageWeight       float32                          `json:"package_weight"`
	PackageHeight       float32                          `json:"package_height"`
	PackageWidth        float32                          `json:"package_width"`
	PackageLength       float32                          `json:"package_length"`
	SKU                 string                           `json:"sku"`
	Name                string                           `json:"name"`
	UPC                 string                           `json:"upc"`
	Quantity            int                              `json:"quantity"`
	OnSale              bool                             `json:"on_sale"`
	Price               float32                          `json:"price"`
	SalePrice           float64                          `json:"sale_price"`
	AttributesAndValues []*ProductAttributeCreationInput `json:"attributes_and_values"`
}

func validateProductUpdateInput(req *http.Request) (*Product, error) {
	product := &Product{}
	decoder := json.NewDecoder(req.Body)
	defer req.Body.Close()
	err := decoder.Decode(product)

	p := structs.New(product)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if p.IsZero() {
		return nil, errors.New("Invalid input provided for product body")
	}

	// we need to be certain that if a user passed us a SKU, that it isn't set
	// to something that mux won't disallow them from retrieving later
	s := p.Field("SKU")
	if !s.IsZero() && !skuValidator.MatchString(product.SKU) {
		return nil, errors.New("Invalid input provided for product SKU")
	}

	product.PackageWeight = float32(Round(float64(product.PackageWeight), .1, 2))
	product.PackageHeight = float32(Round(float64(product.PackageHeight), .1, 2))
	product.PackageWidth = float32(Round(float64(product.PackageWidth), .1, 2))
	product.PackageLength = float32(Round(float64(product.PackageLength), .1, 2))
	product.ProductWeight = float32(Round(float64(product.ProductWeight), .1, 2))
	product.ProductHeight = float32(Round(float64(product.ProductHeight), .1, 2))
	product.ProductWidth = float32(Round(float64(product.ProductWidth), .1, 2))
	product.ProductLength = float32(Round(float64(product.ProductLength), .1, 2))
	product.Price = float32(Round(float64(product.Price), .1, 2))
	product.SalePrice = NullFloat64{sql.NullFloat64{Float64: Round(product.SalePrice.Float64, .1, 2), Valid: true}}

	return product, err
}

func buildProductExistenceHandler(db *sql.DB) http.HandlerFunc {
	// ProductExistenceHandler handles requests to check if a sku exists
	return func(res http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		sku := vars["sku"]

		productExists, err := rowExistsInDB(db, "products", "sku", sku)
		if err != nil {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		responseStatus := http.StatusNotFound
		if productExists {
			responseStatus = http.StatusOK
		}
		res.WriteHeader(responseStatus)
	}
}

// retrieveProductFromDB retrieves a product with a given SKU from the database
func retrieveProductFromDB(db *sql.DB, sku string) (*Product, error) {
	product := &Product{}
	scanArgs := product.generateJoinScanArgs()
	skuJoinRetrievalQuery := buildCompleteProductRetrievalQuery(sku)
	err := db.QueryRow(skuJoinRetrievalQuery, sku).Scan(scanArgs...)
	if err == sql.ErrNoRows {
		return product, errors.Wrap(err, "Error querying for product")
	}

	return product, err
}

func buildSingleProductHandler(db *sql.DB) http.HandlerFunc {
	// SingleProductHandler is a request handler that returns a single Product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		product, err := retrieveProductFromDB(db, sku)
		if err != nil {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		json.NewEncoder(res).Encode(product)
	}
}

func retrieveProductsFromDB(db *sql.DB, queryFilter *QueryFilter) ([]Product, error) {
	var products []Product

	query, args := buildAllProductsRetrievalQuery(queryFilter)
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, "Error encountered querying for products")
	}
	defer rows.Close()
	for rows.Next() {
		var product Product
		_ = rows.Scan(product.generateJoinScanArgs()...)
		products = append(products, product)
	}
	return products, nil
}

func buildProductListHandler(db *sql.DB) http.HandlerFunc {
	// productListHandler is a request handler that returns a list of products
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		products, err := retrieveProductsFromDB(db, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve products from the database")
			return
		}

		productsResponse := &ProductsResponse{
			ListResponse: ListResponse{
				Page:  queryFilter.Page,
				Limit: queryFilter.Limit,
				Count: uint64(len(products)),
			},
			Data: products,
		}
		json.NewEncoder(res).Encode(productsResponse)
	}
}

func deleteProductBySKU(db *sql.DB, sku string) error {
	productDeletionQuery := buildProductDeletionQuery(sku)
	_, err := db.Exec(productDeletionQuery, sku)
	return err
}

func buildProductDeletionHandler(db *sql.DB) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		sku := mux.Vars(req)["sku"]

		// can't delete a product that doesn't exist!
		exists, err := rowExistsInDB(db, "products", "sku", sku)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		err = deleteProductBySKU(db, sku)
		io.WriteString(res, fmt.Sprintf("Successfully deleted product `%s`", sku))
	}
}

func updateProductInDatabase(db *sql.DB, up *Product) error {
	productUpdateQuery, queryArgs := buildProductUpdateQuery(up)
	scanArgs := up.generateScanArgs()
	err := db.QueryRow(productUpdateQuery, queryArgs...).Scan(scanArgs...)
	return err
}

func buildProductUpdateHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// ProductUpdateHandler is a request handler that can update products
		sku := mux.Vars(req)["sku"]

		// can't update a product that doesn't exist!
		exists, err := rowExistsInDB(db, "products", "sku", sku)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "product", sku)
			return
		}

		newerProduct, err := validateProductUpdateInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// eating the error here because we're already certain the sku exists
		existingProduct, err := retrieveProductFromDB(db, sku)
		if err != nil {
			notifyOfInternalIssue(res, err, "merge updated product with existing product")
			return
		}

		// eating the error here because we've already validated input
		mergo.Merge(newerProduct, existingProduct)

		err = updateProductInDatabase(db, newerProduct)
		if err != nil {
			errStr := err.Error()
			notifyOfInternalIssue(res, err, errStr) // "update product in database")
			return
		}

		json.NewEncoder(res).Encode(newerProduct)
	}
}

func validateProductCreationInput(req *http.Request) (*ProductCreationInput, error) {
	pci := &ProductCreationInput{}
	err := json.NewDecoder(req.Body).Decode(pci)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}

	p := structs.New(pci)
	// go will happily decode an invalid input into a completely zeroed struct,
	// so we gotta do checks like this because we're bad at programming.
	if p.IsZero() {
		return nil, errors.New("Invalid input provided for product body")
	}

	// we need to be certain that if a user passed us a SKU, that it isn't set
	// to something that mux won't disallow them from retrieving later
	s := p.Field("SKU")
	if !s.IsZero() && !skuValidator.MatchString(pci.SKU) {
		return nil, errors.New("Invalid input provided for product SKU")
	}

	return pci, err
}

// createProductInDB takes a marshaled Product object and creates an entry for it and a base_product in the database
func createProductInDB(db *sql.DB, np *Product) (int64, error) {
	var newProductID int64
	productCreationQuery, queryArgs := buildProductCreationQuery(np)
	err := db.QueryRow(productCreationQuery, queryArgs...).Scan(&newProductID)
	return newProductID, err
}

func buildProductCreationHandler(db *sql.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		productInput, err := validateProductCreationInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		// can't create a product with a sku that already exists!
		exists, err := rowExistsInDB(db, "products", "sku", productInput.SKU)
		if err != nil || exists {
			notifyOfInvalidRequestBody(res, fmt.Errorf("product with sku `%s` already exists", productInput.SKU))
			return
		}

		progenitor := newProductProgenitorFromProductCreationInput(productInput)
		newProgenitorID, err := createProductProgenitorInDB(db, progenitor)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert product progenitor in database")
			return
		}
		progenitor.ID = newProgenitorID

		for _, attributeAndValues := range productInput.AttributesAndValues {
			_, err = createProductAttributeAndValuesInDBFromInput(db, attributeAndValues, progenitor.ID)
			if err != nil {
				notifyOfInternalIssue(res, err, "insert product attributes and values in database")
				return
			}
		}

		newProduct := &Product{
			ProductProgenitor:   *progenitor,
			ProductProgenitorID: progenitor.ID,
			SKU:                 productInput.SKU,
			Name:                productInput.Name,
			UPC:                 NullString{sql.NullString{String: productInput.UPC, Valid: true}},
			Quantity:            productInput.Quantity,
			Price:               productInput.Price,
			OnSale:              productInput.OnSale,
			SalePrice:           NullFloat64{sql.NullFloat64{Float64: productInput.SalePrice}},
		}

		newProductID, err := createProductInDB(db, newProduct)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert product in database")
			return
		}
		newProduct.ID = newProductID

		json.NewEncoder(res).Encode(newProduct)
	}
}
