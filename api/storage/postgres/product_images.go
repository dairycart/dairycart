package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const assignProductImageIDToProductQuery = `
    UPDATE products
    SET
        primary_image_id = $1,
        updated_on = NOW()
    WHERE id = $2
    RETURNING updated_on;
`

func (pg *postgres) SetPrimaryProductImageForProduct(db storage.Querier, productID, imageID uint64) (t time.Time, err error) {
	err = db.QueryRow(assignProductImageIDToProductQuery, imageID, productID).Scan(&t)
	return t, err
}

const productImageQueryByProductID = `
    SELECT
        id,
        product_root_id,
        thumbnail_url,
        main_url,
        original_url,
        source_url,
        created_on,
        updated_on,
        archived_on
    FROM
        product_images
    WHERE
        archived_on is null
    AND
        product_id = $1
`

func (pg *postgres) GetProductImagesByProductID(db storage.Querier, productID uint64) ([]models.ProductImage, error) {
	var list []models.ProductImage

	rows, err := db.Query(productImageQueryByProductID, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.ProductImage
		err := rows.Scan(
			&p.ID,
			&p.ProductRootID,
			&p.ThumbnailURL,
			&p.MainURL,
			&p.OriginalURL,
			&p.SourceURL,
			&p.CreatedOn,
			&p.UpdatedOn,
			&p.ArchivedOn,
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

const productImageExistenceQuery = `SELECT EXISTS(SELECT id FROM product_images WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) ProductImageExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(productImageExistenceQuery, id).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productImageSelectionQuery = `
    SELECT
        id,
        product_root_id,
        thumbnail_url,
        main_url,
        original_url,
        source_url,
        created_on,
        updated_on,
        archived_on
    FROM
        product_images
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetProductImage(db storage.Querier, id uint64) (*models.ProductImage, error) {
	p := &models.ProductImage{}

	err := db.QueryRow(productImageSelectionQuery, id).Scan(&p.ID, &p.ProductRootID, &p.ThumbnailURL, &p.MainURL, &p.OriginalURL, &p.SourceURL, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

func buildProductImageListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"product_root_id",
			"thumbnail_url",
			"main_url",
			"original_url",
			"source_url",
			"created_on",
			"updated_on",
			"archived_on",
		).
		From("product_images")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetProductImageList(db storage.Querier, qf *models.QueryFilter) ([]models.ProductImage, error) {
	var list []models.ProductImage
	query, args := buildProductImageListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.ProductImage
		err := rows.Scan(
			&p.ID,
			&p.ProductRootID,
			&p.ThumbnailURL,
			&p.MainURL,
			&p.OriginalURL,
			&p.SourceURL,
			&p.CreatedOn,
			&p.UpdatedOn,
			&p.ArchivedOn,
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

func buildProductImageCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("product_images")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetProductImageCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildProductImageCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const productImageCreationQuery = `
    INSERT INTO product_images
        (
            product_root_id, thumbnail_url, main_url, original_url, source_url
        )
    VALUES
        (
            $1, $2, $3, $4, $5
        )
    RETURNING
        id, created_on;
`

func (pg *postgres) CreateProductImage(db storage.Querier, nu *models.ProductImage) (createdID uint64, createdOn time.Time, err error) {
	err = db.QueryRow(productImageCreationQuery, &nu.ProductRootID, &nu.ThumbnailURL, &nu.MainURL, &nu.OriginalURL, &nu.SourceURL).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
}

const productImageUpdateQuery = `
    UPDATE product_images
    SET
        product_root_id = $1,
        thumbnail_url = $2,
        main_url = $3,
        original_url = $4,
        source_url = $5,
        updated_on = NOW()
    WHERE id = $6
    RETURNING updated_on;
`

func (pg *postgres) UpdateProductImage(db storage.Querier, updated *models.ProductImage) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productImageUpdateQuery, &updated.ProductRootID, &updated.ThumbnailURL, &updated.MainURL, &updated.OriginalURL, &updated.SourceURL, &updated.ID).Scan(&t)
	return t, err
}

const productImageDeletionQuery = `
    UPDATE product_images
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteProductImage(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productImageDeletionQuery, id).Scan(&t)
	return t, err
}
