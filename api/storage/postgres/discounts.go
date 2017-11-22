package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"
)

const discountQueryByCode = `
    SELECT
        id,
        name,
        discount_type,
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
    FROM
        discounts
    WHERE
        archived_on is null
    AND
        sku = $1
`

func (pg *postgres) GetDiscountByCode(db storage.Querier, code string) (*models.Discount, error) {
	d := &models.Discount{}

	err := db.QueryRow(discountQueryByCode, code).Scan(&d.ID, &d.Name, &d.DiscountType, &d.Amount, &d.StartsOn, &d.ExpiresOn, &d.RequiresCode, &d.Code, &d.LimitedUse, &d.NumberOfUses, &d.LoginRequired, &d.CreatedOn, &d.UpdatedOn, &d.ArchivedOn)
	return d, err
}

const discountExistenceQuery = `SELECT EXISTS(SELECT id FROM discounts WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) DiscountExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(discountExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const discountSelectionQuery = `
    SELECT
        id,
        name,
        discount_type,
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
    FROM
        discounts
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetDiscount(db storage.Querier, id uint64) (*models.Discount, error) {
	d := &models.Discount{}

	err := db.QueryRow(discountSelectionQuery, id).Scan(&d.ID, &d.Name, &d.DiscountType, &d.Amount, &d.StartsOn, &d.ExpiresOn, &d.RequiresCode, &d.Code, &d.LimitedUse, &d.NumberOfUses, &d.LoginRequired, &d.CreatedOn, &d.UpdatedOn, &d.ArchivedOn)

	return d, err
}

const discountCreationQuery = `
    INSERT INTO discounts
        (
            name, discount_type, amount, starts_on, expires_on, requires_code, code, limited_use, number_of_uses, login_required
        )
    VALUES
        (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateDiscount(db storage.Querier, nu *models.Discount) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)

	err := db.QueryRow(discountCreationQuery, &nu.Name, &nu.DiscountType, &nu.Amount, &nu.StartsOn, &nu.ExpiresOn, &nu.RequiresCode, &nu.Code, &nu.LimitedUse, &nu.NumberOfUses, &nu.LoginRequired).Scan(&createdID, &createdAt)

	return createdID, createdAt, err
}

const discountUpdateQuery = `
    UPDATE discounts
    SET
        name = $1, 
        discount_type = $2, 
        amount = $3, 
        starts_on = $4, 
        expires_on = $5, 
        requires_code = $6, 
        code = $7, 
        limited_use = $8, 
        number_of_uses = $9, 
        login_required = $10, 
        updated_on = NOW()
    WHERE id = $11
    RETURNING updated_on;
`

func (pg *postgres) UpdateDiscount(db storage.Querier, updated *models.Discount) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(discountUpdateQuery, &updated.Name, &updated.DiscountType, &updated.Amount, &updated.StartsOn, &updated.ExpiresOn, &updated.RequiresCode, &updated.Code, &updated.LimitedUse, &updated.NumberOfUses, &updated.LoginRequired, &updated.ID).Scan(&t)
	return t, err
}

const discountDeletionQuery = `
    UPDATE discounts
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteDiscount(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(discountDeletionQuery, id).Scan(&t)
	return t, err
}
