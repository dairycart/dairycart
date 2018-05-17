package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/storage/database"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const productImageBridgeExistenceQuery = `SELECT EXISTS(SELECT id FROM product_image_bridge WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) ProductImageBridgeExists(db database.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(productImageBridgeExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productImageBridgeSelectionQuery = `
    SELECT
        id,
        product_id,
        product_image_id
    FROM
        product_image_bridge
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetProductImageBridge(db database.Querier, id uint64) (*models.ProductImageBridge, error) {
	p := &models.ProductImageBridge{}

	err := db.QueryRow(productImageBridgeSelectionQuery, id).Scan(&p.ID, &p.ProductID, &p.ProductImageID)

	return p, err
}

func buildProductImageBridgeListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"product_id",
			"product_image_id",
		).
		From("product_image_bridge")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetProductImageBridgeList(db database.Querier, qf *models.QueryFilter) ([]models.ProductImageBridge, error) {
	var list []models.ProductImageBridge
	query, args := buildProductImageBridgeListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.ProductImageBridge
		err := rows.Scan(
			&p.ID,
			&p.ProductID,
			&p.ProductImageID,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, err
}

func buildProductImageBridgeCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("product_image_bridge")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetProductImageBridgeCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildProductImageBridgeCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const productImageBridgeCreationQuery = `
    INSERT INTO product_image_bridge
        (
            product_id, product_image_id
        )
    VALUES
        (
            $1, $2
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateProductImageBridge(db database.Querier, nu *models.ProductImageBridge) (createdID uint64, createdOn time.Time, err error) {
	err = db.QueryRow(productImageBridgeCreationQuery, &nu.ProductID, &nu.ProductImageID).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
}

const productImageBridgeUpdateQuery = `
    UPDATE product_image_bridge
    SET
        product_id = $1,
        product_image_id = $2,
        updated_on = NOW()
    WHERE id = $3
    RETURNING updated_on;
`

func (pg *postgres) UpdateProductImageBridge(db database.Querier, updated *models.ProductImageBridge) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productImageBridgeUpdateQuery, &updated.ProductID, &updated.ProductImageID, &updated.ID).Scan(&t)
	return t, err
}

const productImageBridgeDeletionQuery = `
    UPDATE product_image_bridge
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteProductImageBridge(db database.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productImageBridgeDeletionQuery, id).Scan(&t)
	return t, err
}
