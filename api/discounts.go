package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"github.com/imdario/mergo"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	discountDBColumns = `
		id,
		name,
		type,
		amount,
		starts_on,
		expires_on,
		requires_code,
		code,
		limited_use,
		number_of_uses,
		login_required,
		created_on,
		updated_on,
		archived_on
	`

	discountRetrievalQuery = `SELECT * FROM discounts WHERE id = $1`
	discountExistenceQuery = `SELECT EXISTS(SELECT 1 FROM discounts WHERE id = $1 AND archived_on IS NULL)`
	discountDeletionQuery  = `UPDATE discounts SET archived_on = NOW() WHERE id = $1 AND archived_on IS NULL`
)

// Discount represents pricing changes that apply temporarily to products
type Discount struct {
	DBRow
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Amount        float32   `json:"amount"`
	StartsOn      time.Time `json:"starts_on"`
	ExpiresOn     NullTime  `json:"expires_on"`
	RequiresCode  bool      `json:"requires_code"`
	Code          string    `json:"code,omitempty"`
	LimitedUse    bool      `json:"limited_use"`
	NumberOfUses  int64     `json:"number_of_uses,omitempty"`
	LoginRequired bool      `json:"login_required"`
}

// generateScanArgs generates an array of pointers to struct fields for sql.Scan to populate
func (d *Discount) generateScanArgs() []interface{} {
	return []interface{}{
		&d.ID,
		&d.Name,
		&d.Type,
		&d.Amount,
		&d.StartsOn,
		&d.ExpiresOn,
		&d.RequiresCode,
		&d.Code,
		&d.LimitedUse,
		&d.NumberOfUses,
		&d.LoginRequired,
		&d.CreatedOn,
		&d.UpdatedOn,
		&d.ArchivedOn,
	}
}

// DiscountCreationInput represents user input for creating new discounts
type DiscountCreationInput struct {
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Amount        float32   `json:"amount"`
	StartsOn      time.Time `json:"starts_on"`
	ExpiresOn     NullTime  `json:"expires_on"`
	RequiresCode  bool      `json:"requires_code"`
	Code          string    `json:"code"`
	LimitedUse    bool      `json:"limited_use"`
	NumberOfUses  int64     `json:"number_of_uses"`
	LoginRequired bool      `json:"login_required"`
}

// DiscountsResponse is a discount response struct
type DiscountsResponse struct {
	ListResponse
	Data []Discount `json:"data"`
}

func (d *Discount) discountTypeIsValid() bool {
	// Because Go doesn't have typed enums (https://github.com/golang/go/issues/19814),
	// this is my only real line of defense against a user attempting to load an invalid
	// discount type into the database. It's lame, type enums aren't, here's hoping.
	return d.Type == "percentage" || d.Type == "flat_amount"
}

func retrieveDiscountFromDB(db *sqlx.DB, discountID string) (Discount, error) {
	var d Discount
	err := db.Get(&d, discountRetrievalQuery, discountID)
	return d, err
}

func buildDiscountRetrievalHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountRetrievalHandler is a request handler that returns a single Discount
	return func(res http.ResponseWriter, req *http.Request) {
		discountID := mux.Vars(req)["discount_id"]

		discount, err := retrieveDiscountFromDB(db, discountID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "discount", discountID)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving discount from database")
			return
		}

		json.NewEncoder(res).Encode(discount)
	}
}

func buildDiscountListRetrievalHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountListRetrievalHandler is a request handler that returns a list of Discounts
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)
		count, err := getRowCount(db, "discounts", queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of discounts from the database")
			return
		}

		var discounts []Discount
		query, args := buildDiscountListQuery(queryFilter)
		err = retrieveListOfRowsFromDB(db, query, args, &discounts)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve discounts from the database")
			return
		}

		discountsResponse := &DiscountsResponse{
			ListResponse: ListResponse{
				Page:  queryFilter.Page,
				Limit: queryFilter.Limit,
				Count: count,
			},
			Data: discounts,
		}
		json.NewEncoder(res).Encode(discountsResponse)
	}
}

func validateDiscountCreationInput(req *http.Request) (*Discount, error) {
	d := &Discount{}
	err := json.NewDecoder(req.Body).Decode(d)
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	s := structs.New(d)
	if s.IsZero() {
		return nil, errors.New("Invalid input provided for discount body")
	}

	return d, err
}

func createDiscountInDB(db *sqlx.DB, in *Discount) (*Discount, error) {
	discountCreationQuery, queryArgs := buildDiscountCreationQuery(in)
	scanArgs := in.generateScanArgs()
	err := db.QueryRow(discountCreationQuery, queryArgs...).Scan(scanArgs...)
	return in, err
}

func buildDiscountCreationHandler(db *sqlx.DB) http.HandlerFunc {
	// DiscountCreationHandler is a request handler that creates a Discount from user input
	return func(res http.ResponseWriter, req *http.Request) {
		newDiscount, err := validateDiscountCreationInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		newDiscount, err = createDiscountInDB(db, newDiscount)
		if err != nil {
			notifyOfInternalIssue(res, err, "insert discount into database")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newDiscount)
	}
}

func archiveDiscount(db *sqlx.DB, discountID string) error {
	_, err := db.Exec(discountDeletionQuery, discountID)
	return err
}

func buildDiscountDeletionHandler(db *sqlx.DB) http.HandlerFunc {
	// ProductDeletionHandler is a request handler that deletes a single product
	return func(res http.ResponseWriter, req *http.Request) {
		discountID := mux.Vars(req)["discount_id"]

		// can't delete a discount that doesn't exist!
		exists, err := rowExistsInDB(db, discountExistenceQuery, discountID)
		if err != nil || !exists {
			respondThatRowDoesNotExist(req, res, "discount", discountID)
			return
		}

		err = archiveDiscount(db, discountID)
		io.WriteString(res, fmt.Sprintf("Successfully archived discount `%s`", discountID))
	}
}

func validateDiscountUpdateInput(req *http.Request) (*Discount, error) {
	d := &Discount{}

	err := json.NewDecoder(req.Body).Decode(d)
	if err != nil {
		return nil, err
	}

	p := structs.New(d)
	if p.IsZero() {
		return nil, errors.New("Invalid input provided for product body")
	}

	return d, nil
}

func updateDiscountInDatabase(db *sqlx.DB, up *Discount) error {
	discountUpdateQuery, queryArgs := buildDiscountUpdateQuery(up)
	scanArgs := up.generateScanArgs()
	err := db.QueryRow(discountUpdateQuery, queryArgs...).Scan(scanArgs...)
	return err
}

func buildDiscountUpdateHandler(db *sqlx.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// DiscountUpdateHandler is a request handler that can update discounts
		discountID := mux.Vars(req)["discount_id"]

		updatedDiscount, err := validateDiscountUpdateInput(req)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingDiscount, err := retrieveDiscountFromDB(db, discountID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "discount", discountID)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving discount from database")
			return
		}

		// eating the error here because we've already validated input
		mergo.Merge(updatedDiscount, existingDiscount)

		err = updateDiscountInDatabase(db, updatedDiscount)
		if err != nil {
			notifyOfInternalIssue(res, err, "update product in database")
			return
		}

		json.NewEncoder(res).Encode(updatedDiscount)
	}
}
