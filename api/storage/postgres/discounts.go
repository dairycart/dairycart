package postgres

import (
	"time"

	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"
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

func (pg *Postgres) GetDiscountByCode(code string) (models.Discount, error) {
	var d models.Discount

	err := pg.DB.QueryRow(discountQueryByCode, code).Scan(&d.ID, &d.Name, &d.DiscountType, &d.Amount, &d.StartsOn, &d.ExpiresOn, &d.RequiresCode, &d.Code, &d.LimitedUse, &d.NumberOfUses, &d.LoginRequired, &d.CreatedOn, &d.UpdatedOn, &d.ArchivedOn)
	return d, err
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

func (pg *Postgres) GetDiscount(id uint64) (models.Discount, error) {
	var d models.Discount

	err := pg.DB.QueryRow(discountSelectionQuery, id).Scan(&d.ID, &d.Name, &d.DiscountType, &d.Amount, &d.StartsOn, &d.ExpiresOn, &d.RequiresCode, &d.Code, &d.LimitedUse, &d.NumberOfUses, &d.LoginRequired, &d.CreatedOn, &d.UpdatedOn, &d.ArchivedOn)

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

func (pg *Postgres) CreateDiscount(np models.Discount) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)
	err := pg.DB.QueryRow(discountCreationQuery, &np.Name, &np.DiscountType, &np.Amount, &np.StartsOn, &np.ExpiresOn, &np.RequiresCode, &np.Code, &np.LimitedUse, &np.NumberOfUses, &np.LoginRequired).Scan(&createdID, &createdAt)

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

func (pg *Postgres) UpdateDiscount(updated models.Discount) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(discountUpdateQuery, &updated.Name, &updated.DiscountType, &updated.Amount, &updated.StartsOn, &updated.ExpiresOn, &updated.RequiresCode, &updated.Code, &updated.LimitedUse, &updated.NumberOfUses, &updated.LoginRequired, &updated.ID).Scan(&t)
	return t, err
}

const discountDeletionQuery = `
    UPDATE discounts
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *Postgres) DeleteDiscount(id uint64) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(discountDeletionQuery, id).Scan(&t)
	return t, err
}
