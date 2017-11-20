package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

const productVariantBridgeExistenceQuery = `SELECT EXISTS(SELECT id FROM product_variant_bridge WHERE id = $1 and archived_on IS NULL);`

func (pg *Postgres) ProductVariantBridgeExists(id uint64) (bool, error) {
	var exists string

	err := pg.DB.QueryRow(productVariantBridgeExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productVariantBridgeSelectionQuery = `
    SELECT
        id,
        product_id,
        product_option_value_id,
        created_on,
        archived_on
    FROM
        product_variant_bridge
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *Postgres) GetProductVariantBridge(id uint64) (*models.ProductVariantBridge, error) {
	p := &models.ProductVariantBridge{}

	err := pg.DB.QueryRow(productVariantBridgeSelectionQuery, id).Scan(&p.ID, &p.ProductID, &p.ProductOptionValueID, &p.CreatedOn, &p.ArchivedOn)

	return p, err
}

const productvariantbridgeCreationQuery = `
    INSERT INTO product_variant_bridge
        (
            product_id, product_option_value_id
        )
    VALUES
        (
            $1, $2
        )
    RETURNING
        id, created_on;
`

func (pg *Postgres) CreateProductVariantBridge(nu *models.ProductVariantBridge) (uint64, time.Time, error) {
	var (
		createdID uint64
		createdAt time.Time
	)
	err := pg.DB.QueryRow(productvariantbridgeCreationQuery, &nu.ProductID, &nu.ProductOptionValueID).Scan(&createdID, &createdAt)

	return createdID, createdAt, err
}

const productVariantBridgeUpdateQuery = `
    UPDATE product_variant_bridge
    SET
        product_id = $1, 
        product_option_value_id = $2
    WHERE id = $2
    RETURNING updated_on;
`

func (pg *Postgres) UpdateProductVariantBridge(updated *models.ProductVariantBridge) (time.Time, error) {
	var t time.Time
	err := pg.DB.QueryRow(productVariantBridgeUpdateQuery, &updated.ProductID, &updated.ProductOptionValueID, &updated.ID).Scan(&t)
	return t, err
}

const productVariantBridgeDeletionQuery = `
    UPDATE product_variant_bridge
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *Postgres) DeleteProductVariantBridge(id uint64, tx *sql.Tx) (t time.Time, err error) {
	if tx != nil {
		err = tx.QueryRow(productVariantBridgeDeletionQuery, id).Scan(&t)
	} else {
		err = pg.DB.QueryRow(productVariantBridgeDeletionQuery, id).Scan(&t)
	}
	return
}
