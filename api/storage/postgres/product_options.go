package postgres

import (
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

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

func (pg *Postgres) GetProductOption(id uint64) (*models.ProductOption, error) {
	p := &models.ProductOption{}

	err := pg.DB.QueryRow(productOptionSelectionQuery, id).Scan(&p.ID, &p.Name, &p.ProductRootID, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

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

func (pg *Postgres) CreateProductOption(nu *models.ProductOption) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)
	err := pg.DB.QueryRow(productoptionCreationQuery, &nu.Name, &nu.ProductRootID).Scan(&createdID, &createdAt)

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

func (pg *Postgres) UpdateProductOption(updated *models.ProductOption) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(productOptionUpdateQuery, &updated.Name, &updated.ProductRootID, &updated.ID).Scan(&t)
	return t, err
}

const productOptionDeletionQuery = `
    UPDATE product_options
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *Postgres) DeleteProductOption(id uint64) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(productOptionDeletionQuery, id).Scan(&t)
	return t, err
}
