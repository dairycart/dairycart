package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

const productOptionValueExistenceQuery = `SELECT EXISTS(SELECT id FROM product_option_values WHERE id = $1 and archived_on IS NULL);`

func (pg *Postgres) ProductOptionValueExists(id uint64) (bool, error) {
	var exists string

	err := pg.DB.QueryRow(productOptionValueExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productOptionValueSelectionQuery = `
    SELECT
        id,
        product_option_id,
        value,
        created_on,
        updated_on,
        archived_on
    FROM
        product_option_values
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *Postgres) GetProductOptionValue(id uint64) (*models.ProductOptionValue, error) {
	p := &models.ProductOptionValue{}

	err := pg.DB.QueryRow(productOptionValueSelectionQuery, id).Scan(&p.ID, &p.ProductOptionID, &p.Value, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

const productoptionvalueCreationQuery = `
    INSERT INTO product_option_values
        (
            product_option_id, value
        )
    VALUES
        (
            $1, $2
        )
    RETURNING
        id, created_on;
`

func (pg *Postgres) CreateProductOptionValue(nu *models.ProductOptionValue) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)
	err := pg.DB.QueryRow(productoptionvalueCreationQuery, &nu.ProductOptionID, &nu.Value).Scan(&createdID, &createdAt)

	return createdID, createdAt, err
}

const productOptionValueUpdateQuery = `
    UPDATE product_option_values
    SET
        product_option_id = $1, 
        value = $2, 
        updated_on = NOW()
    WHERE id = $3
    RETURNING updated_on;
`

func (pg *Postgres) UpdateProductOptionValue(updated *models.ProductOptionValue) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(productOptionValueUpdateQuery, &updated.ProductOptionID, &updated.Value, &updated.ID).Scan(&t)
	return t, err
}

const productOptionValueDeletionQuery = `
    UPDATE product_option_values
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *Postgres) DeleteProductOptionValue(id uint64, tx *sql.Tx) (t time.Time, err error) {
	if tx != nil {
		err = tx.QueryRow(productOptionValueDeletionQuery, id).Scan(&t)
	} else {
		err = pg.DB.QueryRow(productOptionValueDeletionQuery, id).Scan(&t)
	}
	return
}
