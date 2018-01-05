package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const productOptionValueForOptionIDExistenceQuery = `SELECT EXISTS(SELECT id FROM product_option_values WHERE product_option_id = $1 AND value = $2 and archived_on IS NULL);`

func (pg *postgres) ProductOptionValueForOptionIDExists(db storage.Querier, optionID uint64, value string) (bool, error) {
	var exists string

	err := db.QueryRow(productOptionValueForOptionIDExistenceQuery, optionID, value).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return exists == "true", err
}

const productOptionValueArchiveQueryByOptionID = `
    UPDATE product_option_values
    SET archived_on = NOW()
    WHERE product_option_id = $1
    RETURNING archived_on
`

func (pg *postgres) ArchiveProductOptionValuesForOption(db storage.Querier, optionID uint64) (t time.Time, err error) {
	err = db.QueryRow(productOptionValueArchiveQueryByOptionID, optionID).Scan(&t)
	return t, err
}

const productOptionValueRetrievalQueryByOptionID = `
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
        product_option_id = $1
`

func (pg *postgres) GetProductOptionValuesForOption(db storage.Querier, optionID uint64) ([]models.ProductOptionValue, error) {
	var list []models.ProductOptionValue

	rows, err := db.Query(productOptionValueRetrievalQueryByOptionID, optionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.ProductOptionValue
		err := rows.Scan(
			&p.ID,
			&p.ProductOptionID,
			&p.Value,
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

const productOptionValueExistenceQuery = `SELECT EXISTS(SELECT id FROM product_option_values WHERE id = $1 and archived_on IS NULL);`

func (pg *postgres) ProductOptionValueExists(db storage.Querier, id uint64) (bool, error) {
	var exists string

	err := db.QueryRow(productOptionValueExistenceQuery, id).Scan(&exists)
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

func (pg *postgres) GetProductOptionValue(db storage.Querier, id uint64) (*models.ProductOptionValue, error) {
	p := &models.ProductOptionValue{}

	err := db.QueryRow(productOptionValueSelectionQuery, id).Scan(&p.ID, &p.ProductOptionID, &p.Value, &p.CreatedOn, &p.UpdatedOn, &p.ArchivedOn)

	return p, err
}

func buildProductOptionValueListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"product_option_id",
			"value",
			"created_on",
			"updated_on",
			"archived_on",
		).
		From("product_option_values")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetProductOptionValueList(db storage.Querier, qf *models.QueryFilter) ([]models.ProductOptionValue, error) {
	var list []models.ProductOptionValue
	query, args := buildProductOptionValueListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p models.ProductOptionValue
		err := rows.Scan(
			&p.ID,
			&p.ProductOptionID,
			&p.Value,
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

func buildProductOptionValueCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("product_option_values")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetProductOptionValueCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildProductOptionValueCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const productOptionValueCreationQuery = `
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

func (pg *postgres) CreateProductOptionValue(db storage.Querier, nu *models.ProductOptionValue) (createdID uint64, createdOn time.Time, err error) {
	err = db.QueryRow(productOptionValueCreationQuery, &nu.ProductOptionID, &nu.Value).Scan(&createdID, &createdOn)
	return createdID, createdOn, err
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

func (pg *postgres) UpdateProductOptionValue(db storage.Querier, updated *models.ProductOptionValue) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(productOptionValueUpdateQuery, &updated.ProductOptionID, &updated.Value, &updated.ID).Scan(&t)
	return t, err
}

const productOptionValueDeletionQuery = `
    UPDATE product_option_values
    SET archived_on = NOW()
    WHERE id = $1
    RETURNING archived_on
`

func (pg *postgres) DeleteProductOptionValue(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productOptionValueDeletionQuery, id).Scan(&t)
	return t, err
}

const productOptionValueWithProductRootIDDeletionQuery = `
    UPDATE product_option_values
	SET archived_on = NOW()
	WHERE product_option_id IN (SELECT id FROM product_options WHERE product_root_id = $1)
`

func (pg *postgres) ArchiveProductOptionValuesWithProductRootID(db storage.Querier, id uint64) (t time.Time, err error) {
	err = db.QueryRow(productOptionValueWithProductRootIDDeletionQuery, id).Scan(&t)
	return t, err
}
