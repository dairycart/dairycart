package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"
)

const productOptionExistenceQuery = `SELECT EXISTS(SELECT id FROM product_options WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) ProductOptionExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(productOptionExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productOptionSelectionQuery = `
    SELECT
        id,
        name,
        product_root_id,
        created_on,
        updated_on,
        archived_on
    FROM
        product_options
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetProductOption(db storage.Querier, id uint64) (*models.ProductOption, error) {
	p := &models.ProductOption{}

	err := db.QueryRow(productOptionSelectionQuery, id).Scan(&p.ID, &p.Name, &p.ProductRootID, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

const productoptionCreationQuery = `
    INSERT INTO product_options
        (
            name, product_root_id
        )
    VALUES
        (
            $1, $2
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateProductOption(db storage.Querier, nu *models.ProductOption) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)
	err := db.QueryRow(productoptionCreationQuery, &nu.Name, &nu.ProductRootID).Scan(&createdID, &createdAt)

	return createdID, createdAt, err
}

const productOptionUpdateQuery = `
    UPDATE product_options
    SET
        name = $1, 
        product_root_id = $2, 
        updated_on = NOW()
    WHERE id = $3
    RETURNING updated_on;
`

func (pg *postgres) UpdateProductOption(db storage.Querier, updated *models.ProductOption) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productOptionUpdateQuery, &updated.Name, &updated.ProductRootID, &updated.ID).Scan(&t)
	return t, err
}

const productOptionDeletionQuery = `
    UPDATE product_options
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteProductOption(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productOptionDeletionQuery, id).Scan(&t)
	return t, err
}
