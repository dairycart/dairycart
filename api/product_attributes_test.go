package api

import (
	"database/sql/driver"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var exampleProductAttribute *ProductAttribute
var productAttributeHeaders []string
var productAttributeData []driver.Value

func init() {
	exampleProductAttribute = &ProductAttribute{
		ID:   123,
		Name: "attribute",
	}
	productAttributeHeaders = []string{"id", "name", "product_progenitor_id", "created_at", "updated_at", "archived_at"}
	productAttributeData = []driver.Value{1, 2, "Attribute", exampleTime, nil, nil}
}

func setExpectationsForProductAttributeExistence(mock sqlmock.Sqlmock, id int64, exists bool, err error) {
	exampleRows := sqlmock.NewRows([]string{""}).AddRow(strconv.FormatBool(exists))
	query := buildProductAttributeexistenceQuery(id)
	stringID := strconv.Itoa(int(id))
	mock.ExpectQuery(formatQueryForSQLMock(query)).
		WithArgs(stringID).
		WillReturnRows(exampleRows).
		WillReturnError(err)
}

//////////////////////////////////////////////////////
//                                                  //
//   HTTP                                 :e        //
//     Handler                           'M$\       //
//         Tests                        sf$$br      //
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

func setExpectationsForProductAttributeListQuery(mock sqlmock.Sqlmock, err error) {
	exampleRows := sqlmock.NewRows(productAttributeHeaders).
		AddRow(productAttributeData...).
		AddRow(productAttributeData...).
		AddRow(productAttributeData...)
	query := buildProductAttributeListQuery(strconv.Itoa(int(exampleProgenitor.ID)), defaultQueryFilter)
	mock.ExpectQuery(formatQueryForSQLMock(query)).WillReturnRows(exampleRows).WillReturnError(err)
}

func TestProductAttributeListHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeListQuery(mock, nil)

	productAttributeEndpoint := buildRoute("product_attributes", strconv.Itoa(int(exampleProgenitor.ID)))
	req, err := http.NewRequest("GET", productAttributeEndpoint, nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 200, "status code should be 200")

	expected := &ProductsResponse{
		ListResponse: ListResponse{
			Page:  1,
			Limit: 25,
			Count: 3,
		},
	}

	actual := &ProductsResponse{}
	err = json.NewDecoder(strings.NewReader(res.Body.String())).Decode(actual)
	assert.Nil(t, err)

	assert.Equal(t, expected.Page, actual.Page, "expected and actual product pages should be equal")
	assert.Equal(t, expected.Limit, actual.Limit, "expected and actual product limits should be equal")
	assert.Equal(t, expected.Count, actual.Count, "expected and actual product counts should be equal")
	assert.Equal(t, uint64(len(actual.Data)), actual.Count, "actual product counts and product response count field should be equal")
	ensureExpectationsWereMet(t, mock)
}

func TestProductAttributeListHandlerWithDBErrors(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)
	defer db.Close()
	res, router := setupMockRequestsAndMux(db)
	setExpectationsForProductAttributeListQuery(mock, arbitraryError)

	productAttributeEndpoint := buildRoute("product_attributes", strconv.Itoa(int(exampleProgenitor.ID)))
	req, err := http.NewRequest("GET", productAttributeEndpoint, nil)
	assert.Nil(t, err)

	router.ServeHTTP(res, req)
	assert.Equal(t, res.Code, 500, "status code should be 500")
	ensureExpectationsWereMet(t, mock)
}
