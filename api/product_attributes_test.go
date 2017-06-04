package api

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

const (
	exampleProductAttributeCreationBody = `
		{
			"name": "something",
			"values": [
				"one",
				"two",
				"three"
			]
		}
	`
	exampleProductAttributeUpdateBody = `
		{
			"name": "something else"
		}
	`
)

var exampleProductAttribute *ProductAttribute
var exampleUpdatedProductAttribute *ProductAttribute
var expectedCreatedProductAttribute *ProductAttribute
var exampleProductAttributeInput *ProductAttributeCreationInput
var productAttributeHeaders []string
var productAttributeData []driver.Value

func init() {
	exampleProductAttribute = &ProductAttribute{
		ID:                  123,
		Name:                "something",
		ProductProgenitorID: 2, // == exampleProgenitor.ID
		CreatedAt:           exampleTime,
	}
	exampleUpdatedProductAttribute = &ProductAttribute{
		ID:                  exampleProductAttribute.ID,
		Name:                "something else",
		ProductProgenitorID: exampleProductAttribute.ProductProgenitorID,
	}
	productAttributeHeaders = []string{"id", "name", "product_progenitor_id", "created_at", "updated_at", "archived_at"}

	expectedCreatedProductAttribute = &ProductAttribute{
		ID:                  exampleProductAttribute.ID,
		Name:                "something",
		ProductProgenitorID: exampleProductAttribute.ProductProgenitorID,
		Values: []*ProductAttributeValue{
			{
				ID:                 256, // == exampleProductAttributeValue.ID,
				ProductAttributeID: exampleProductAttribute.ID,
				Value:              "one",
			}, {
				ID:                 256, // == exampleProductAttributeValue.ID,
				ProductAttributeID: exampleProductAttribute.ID,
				Value:              "two",
			}, {
				ID:                 256, // == exampleProductAttributeValue.ID,
				ProductAttributeID: exampleProductAttribute.ID,
				Value:              "three",
			},
		},
	}

	exampleProductAttributeInput = &ProductAttributeCreationInput{
		Name:   "something",
		Values: []string{"one", "two", "three"},
	}
}

func setExpectationsForProductAttributeExistenceByID(mock sqlmock.Sqlmock, a *ProductAttribute, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := buildProductAttributeExistenceQuery(a.ID)
	stringID := strconv.Itoa(int(a.ID))
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeExistenceByName(mock sqlmock.Sqlmock, a *ProductAttribute, progenitorID string, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := buildProductAttributeExistenceQueryForProductByName(a.Name, progenitorID)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(a.Name, progenitorID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeRetrievalQuery(mock sqlmock.Sqlmock, a *ProductAttribute, err error) {
	exampleRows := sqlmock.NewRows(productAttributeHeaders).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...)
	query := formatQueryForSQLMock(buildProductAttributeRetrievalQuery(a.ID))
	mock.ExpectQuery(query).
		WithArgs(a.ID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeListQuery(mock sqlmock.Sqlmock, a *ProductAttribute, err error) {
	exampleRows := sqlmock.NewRows(productAttributeHeaders).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...)
	query := buildProductAttributeListQuery(strconv.Itoa(int(exampleProgenitor.ID)), defaultQueryFilter)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeListQueryWithCount(mock sqlmock.Sqlmock, a *ProductAttribute, err error) {
	exampleRows := sqlmock.NewRows([]string{"count", "id", "name", "product_progenitor_id", "created_at", "updated_at", "archived_at"}).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...).
		AddRow([]driver.Value{3, a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...)
	query := buildProductAttributeListQueryWithCount(strconv.Itoa(int(exampleProgenitor.ID)), defaultQueryFilter)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeCreation(mock sqlmock.Sqlmock, a *ProductAttribute, err error) {
	exampleRows := sqlmock.NewRows([]string{"id"}).AddRow(exampleProductAttribute.ID)
	query, args := buildProductAttributeCreationQuery(a)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func setExpectationsForProductAttributeUpdate(mock sqlmock.Sqlmock, a *ProductAttribute, err error) {
	exampleRows := sqlmock.NewRows(productAttributeHeaders).
		AddRow([]driver.Value{a.ID, a.Name, a.ProductProgenitorID, exampleTime, nil, nil}...)
	query, args := buildProductAttributeUpdateQuery(a)
	queryArgs := argsToDriverValues(args)
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(queryArgs...).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

func TestValidateProductAttributeUpdateInput(t *testing.T) {
	t.Parallel()
	expected := &ProductAttributeUpdateInput{Name: "something else"}
	exampleInput := strings.NewReader(exampleProductAttributeUpdateBody)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	actual, err := validateProductAttributeUpdateInput(req)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "ProductAttributeUpdateInput should match expectation")
}

func TestValidateProductAttributeUpdateInputWithInvalidInput(t *testing.T) {
	t.Parallel()
	exampleInput := strings.NewReader(exampleGarbageInput)

	req := httptest.NewRequest("GET", "http://example.com", exampleInput)
	_, err := validateProductAttributeUpdateInput(req)

	assert.NotNil(t, err)
}

func TestRetrieveProductAttributeFromDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeRetrievalQuery(mock, exampleProductAttribute, nil)

	actual, err := retrieveProductAttributeFromDB(db, exampleProductAttribute.ID)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductAttribute, actual, "expected and actual product attributes should match")
	ensureExpectationsWereMet(t, mock)
}

func TestRetrieveProductAttributeFromDBWithNoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeRetrievalQuery(mock, exampleProductAttribute, sql.ErrNoRows)

	_, err = retrieveProductAttributeFromDB(db, exampleProductAttribute.ID)
	assert.NotNil(t, err)
	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductAttributeInDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductAttributeCreation(mock, exampleProductAttribute, nil)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	actual, err := createProductAttributeInDB(tx, exampleProductAttribute)
	assert.Nil(t, err)
	assert.Equal(t, exampleProductAttribute, actual, "Creating a product attribute should return the created product attribute")

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductAttributeAndValuesInDBFromInput(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductAttributeCreation(mock, expectedCreatedProductAttribute, nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[0], nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[1], nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[2], nil)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	actual, err := createProductAttributeAndValuesInDBFromInput(tx, exampleProductAttributeInput, expectedCreatedProductAttribute.ProductProgenitorID)
	assert.Nil(t, err)
	assert.Equal(t, expectedCreatedProductAttribute, actual, "output from product attribute creation should match expectation")

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductAttributeAndValuesInDBFromInputWithIssueCreatingAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductAttributeCreation(mock, expectedCreatedProductAttribute, arbitraryError)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	_, err = createProductAttributeAndValuesInDBFromInput(tx, exampleProductAttributeInput, expectedCreatedProductAttribute.ProductProgenitorID)
	assert.NotNil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, mock)
}

func TestCreateProductAttributeAndValuesInDBFromInputWithIssueCreatingAttributeValue(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	mock.ExpectBegin()
	setExpectationsForProductAttributeCreation(mock, expectedCreatedProductAttribute, nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[0], arbitraryError)
	mock.ExpectCommit()

	tx, err := db.Begin()
	assert.Nil(t, err)

	_, err = createProductAttributeAndValuesInDBFromInput(tx, exampleProductAttributeInput, expectedCreatedProductAttribute.ProductProgenitorID)
	assert.NotNil(t, err)

	err = tx.Commit()
	assert.Nil(t, err)

	ensureExpectationsWereMet(t, mock)
}

func TestUpdateProductAttributeInDB(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	setExpectationsForProductAttributeUpdate(mock, expectedCreatedProductAttribute, nil)

	err = updateProductAttributeInDB(db, exampleProductAttribute)
	assert.Nil(t, err)
	ensureExpectationsWereMet(t, mock)
}

//////////////////////////////////////////////////////
//                                                  //
//     HTTP                               :e        //
//       Handler                         'M$\       //
//           Tests                      sf$$br      //
//                                    J\J\J$L$L     //
//                                  :d  )fM$$$$$r   //
//                             ..P*\ .4MJP   '*\    //
//                    sed"""""" ser d$$$F           //
//                .M\  ..JM$$$B$$$$BJ$MR  ...       //
//               dF  nMMM$$$R$$$$$$$h"$ks$$"$$r     //
//             J\.. .MMM8$$$$$LM$P\..'**\    *\     //
//            d :d$r "M$$$$br'$M\d$R                //
//           J\MM\ *L   *M$B8MM$B.**                //
//          :fd$>  :fhr 'MRM$$M$$"                  //
//          MJ$>    '5J5..M8$$>                     //
//         :fMM     d$Fd$$R$$F                      //
//         4M$P .$$*.J*$$**                         //
//         M4$> '$>dRdF                             //
//         MMM\   *L*B.                             //
//        :$$F     ?k"Re                            //
//      .$$P\        **'$$B...                      //
//   :e$F"               '""""                      //
//                                                  //
//////////////////////////////////////////////////////

func TestProductAttributeListHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeListQueryWithCount(mock, exampleProductAttribute, nil)
	setExpectationsForProductAttributeValueRetrievalByAttributeID(mock, exampleProductAttribute, nil)
	setExpectationsForProductAttributeValueRetrievalByAttributeID(mock, exampleProductAttribute, nil)
	setExpectationsForProductAttributeValueRetrievalByAttributeID(mock, exampleProductAttribute, nil)

	productAttributeEndpoint := buildRoute("product_attributes", strconv.Itoa(int(exampleProgenitor.ID)))
	req, err := http.NewRequest("GET", productAttributeEndpoint, nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Code, "status code should be 200")

	expected := &ProductAttributesResponse{
		ListResponse: ListResponse{
			Page:  1,
			Limit: 25,
			Count: 3,
		},
	}

	actual := &ProductAttributesResponse{}
	bodyString := res.Body.String()
	err = json.NewDecoder(strings.NewReader(bodyString)).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product attribute pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product attribute limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product attribute counts should be equal")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeListHandlerWithErrorsRetrievingValues(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeListQueryWithCount(mock, exampleProductAttribute, nil)
	setExpectationsForProductAttributeValueRetrievalByAttributeID(mock, exampleProductAttribute, nil)
	setExpectationsForProductAttributeValueRetrievalByAttributeID(mock, exampleProductAttribute, nil)
	setExpectationsForProductAttributeValueRetrievalByAttributeID(mock, exampleProductAttribute, arbitraryError)

	productAttributeEndpoint := buildRoute("product_attributes", strconv.Itoa(int(exampleProgenitor.ID)))
	req, err := http.NewRequest("GET", productAttributeEndpoint, nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeListHandlerWithDBErrors(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeListQueryWithCount(mock, exampleProductAttribute, arbitraryError)

	productAttributeEndpoint := buildRoute("product_attributes", strconv.Itoa(int(exampleProgenitor.ID)))
	req, err := http.NewRequest("GET", productAttributeEndpoint, nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeCreationHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductAttribute.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductAttributeExistenceByName(mock, expectedCreatedProductAttribute, progenitorIDString, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductAttributeCreation(mock, expectedCreatedProductAttribute, nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[0], nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[1], nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[2], nil)
	mock.ExpectCommit()

	productAttributeEndpoint := buildRoute("product_attributes", progenitorIDString)
	req, err := http.NewRequest("POST", productAttributeEndpoint, strings.NewReader(exampleProductAttributeCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeCreationHandlerFailureToSetupTransaction(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductAttribute.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductAttributeExistenceByName(mock, expectedCreatedProductAttribute, progenitorIDString, false, nil)
	mock.ExpectBegin().WillReturnError(arbitraryError)

	productAttributeEndpoint := buildRoute("product_attributes", progenitorIDString)
	req, err := http.NewRequest("POST", productAttributeEndpoint, strings.NewReader(exampleProductAttributeCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeCreationHandlerFailureToCommitTransaction(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductAttribute.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductAttributeExistenceByName(mock, expectedCreatedProductAttribute, progenitorIDString, false, nil)
	mock.ExpectBegin()
	setExpectationsForProductAttributeCreation(mock, expectedCreatedProductAttribute, nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[0], nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[1], nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[2], nil)
	mock.ExpectCommit().WillReturnError(arbitraryError)

	productAttributeEndpoint := buildRoute("product_attributes", progenitorIDString)
	req, err := http.NewRequest("POST", productAttributeEndpoint, strings.NewReader(exampleProductAttributeCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeCreationHandlerWhenAttributeWithTheSameNameCheckReturnsNoRows(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductAttribute.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductAttributeExistenceByName(mock, expectedCreatedProductAttribute, progenitorIDString, false, sql.ErrNoRows)
	mock.ExpectBegin()
	setExpectationsForProductAttributeCreation(mock, expectedCreatedProductAttribute, nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[0], nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[1], nil)
	setExpectationsForProductAttributeValueCreation(mock, expectedCreatedProductAttribute.Values[2], nil)
	mock.ExpectCommit()

	productAttributeEndpoint := buildRoute("product_attributes", progenitorIDString)
	req, err := http.NewRequest("POST", productAttributeEndpoint, strings.NewReader(exampleProductAttributeCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeCreationHandlerWithNonExistentProgenitor(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductAttribute.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, false)

	productAttributeEndpoint := buildRoute("product_attributes", progenitorIDString)
	req, err := http.NewRequest("POST", productAttributeEndpoint, strings.NewReader(exampleProductAttributeCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeCreationHandlerWithInvalidAttributeCreationInput(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductAttribute.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)

	productAttributeEndpoint := buildRoute("product_attributes", progenitorIDString)
	req, err := http.NewRequest("POST", productAttributeEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeCreationHandlerWhenAttributeWithTheSameNameExists(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductAttribute.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductAttributeExistenceByName(mock, expectedCreatedProductAttribute, progenitorIDString, true, nil)

	productAttributeEndpoint := buildRoute("product_attributes", progenitorIDString)
	req, err := http.NewRequest("POST", productAttributeEndpoint, strings.NewReader(exampleProductAttributeCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeCreationHandlerWithProblemsCreatingAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	progenitorIDString := strconv.Itoa(int(exampleProductAttribute.ProductProgenitorID))
	setExpectationsForProductProgenitorExistence(mock, progenitorIDString, true)
	setExpectationsForProductAttributeExistenceByName(mock, expectedCreatedProductAttribute, progenitorIDString, false, nil)

	mock.ExpectBegin()
	setExpectationsForProductAttributeCreation(mock, expectedCreatedProductAttribute, arbitraryError)
	mock.ExpectRollback()

	productAttributeEndpoint := buildRoute("product_attributes", progenitorIDString)
	req, err := http.NewRequest("POST", productAttributeEndpoint, strings.NewReader(exampleProductAttributeCreationBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeUpdateHandler(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttribute.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeRetrievalQuery(mock, exampleProductAttribute, nil)
	setExpectationsForProductAttributeUpdate(mock, exampleUpdatedProductAttribute, nil)

	productAttributeEndpoint := buildRoute("product_attributes", attributeIDString)
	req, err := http.NewRequest("PUT", productAttributeEndpoint, strings.NewReader(exampleProductAttributeUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code, "status code should be 200")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeUpdateHandlerWithNonexistentAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttribute.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, false, nil)

	productAttributeEndpoint := buildRoute("product_attributes", attributeIDString)
	req, err := http.NewRequest("PUT", productAttributeEndpoint, strings.NewReader(exampleProductAttributeUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code, "status code should be 404")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeUpdateHandlerWithInvalidInput(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttribute.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)

	productAttributeEndpoint := buildRoute("product_attributes", attributeIDString)
	req, err := http.NewRequest("PUT", productAttributeEndpoint, strings.NewReader(exampleGarbageInput))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code, "status code should be 400")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeUpdateHandlerWithErrorRetrievingAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttribute.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeRetrievalQuery(mock, exampleProductAttribute, arbitraryError)

	productAttributeEndpoint := buildRoute("product_attributes", attributeIDString)
	req, err := http.NewRequest("PUT", productAttributeEndpoint, strings.NewReader(exampleProductAttributeUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeUpdateHandlerWithErrorUpdatingAttribute(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	attributeIDString := strconv.Itoa(int(exampleProductAttribute.ID))

	setExpectationsForProductAttributeExistenceByID(mock, exampleProductAttribute, true, nil)
	setExpectationsForProductAttributeRetrievalQuery(mock, exampleProductAttribute, nil)
	setExpectationsForProductAttributeUpdate(mock, exampleUpdatedProductAttribute, arbitraryError)

	productAttributeEndpoint := buildRoute("product_attributes", attributeIDString)
	req, err := http.NewRequest("PUT", productAttributeEndpoint, strings.NewReader(exampleProductAttributeUpdateBody))
	assert.Nil(t, err)
	router.ServeHTTP(res, req)

	assert.Equal(t, 500, res.Code, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
