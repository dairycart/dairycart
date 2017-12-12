package postgres

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

const discountQueryByCode = `
    SELECT
        id,
        amount,
        limited_use,
        requires_code,
        discount_type,
        archived_on,
        code,
        expires_on,
        updated_on,
        starts_on,
        created_on,
        number_of_uses,
        name,
        login_required
    FROM
        discounts
    WHERE
        archived_on is null
    AND
        sku = $1
`

func (pg *postgres) GetDiscountByCode(db storage.Querier, code string) (*models.Discount, error) {
	d := &models.Discount{}
	err := db.QueryRow(discountQueryByCode, code).Scan(&d.ID, &d.Amount, &d.LimitedUse, &d.RequiresCode, &d.DiscountType, &d.ArchivedOn, &d.Code, &d.ExpiresOn, &d.UpdatedOn, &d.StartsOn, &d.CreatedOn, &d.NumberOfUses, &d.Name, &d.LoginRequired)
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
        amount,
        limited_use,
        requires_code,
        discount_type,
        archived_on,
        code,
        expires_on,
        updated_on,
        starts_on,
        created_on,
        number_of_uses,
        name,
        login_required
    FROM
        discounts
    WHERE
        archived_on is null
    AND
        id = $1
`

func (pg *postgres) GetDiscount(db storage.Querier, id uint64) (*models.Discount, error) {
	d := &models.Discount{}

	err := db.QueryRow(discountSelectionQuery, id).Scan(&d.ID, &d.Amount, &d.LimitedUse, &d.RequiresCode, &d.DiscountType, &d.ArchivedOn, &d.Code, &d.ExpiresOn, &d.UpdatedOn, &d.StartsOn, &d.CreatedOn, &d.NumberOfUses, &d.Name, &d.LoginRequired)

	return d, err
}

func buildDiscountListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"amount",
			"limited_use",
			"requires_code",
			"discount_type",
			"archived_on",
			"code",
			"expires_on",
			"updated_on",
			"starts_on",
			"created_on",
			"number_of_uses",
			"name",
			"login_required",
		).
		From("discounts")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, true).ToSql()
	return query, args
}

func (pg *postgres) GetDiscountList(db storage.Querier, qf *models.QueryFilter) ([]models.Discount, error) {
	var list []models.Discount
	query, args := buildDiscountListRetrievalQuery(qf)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var d models.Discount
		err := rows.Scan(
			&d.ID,
			&d.Amount,
			&d.LimitedUse,
			&d.RequiresCode,
			&d.DiscountType,
			&d.ArchivedOn,
			&d.Code,
			&d.ExpiresOn,
			&d.UpdatedOn,
			&d.StartsOn,
			&d.CreatedOn,
			&d.NumberOfUses,
			&d.Name,
			&d.LoginRequired,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, d)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return list, err
}

func buildDiscountCountRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	queryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(id)").
		From("discounts")

	query, args, _ := applyQueryFilterToQueryBuilder(queryBuilder, qf, false).ToSql()
	return query, args
}

func (pg *postgres) GetDiscountCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	var count uint64
	query, args := buildDiscountCountRetrievalQuery(qf)
	err := db.QueryRow(query, args...).Scan(&count)
	return count, err
}

const discountCreationQuery = `
    INSERT INTO discounts
        (
            amount, limited_use, requires_code, discount_type, code, expires_on, starts_on, number_of_uses, name, login_required
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

	err := db.QueryRow(discountCreationQuery, &nu.Amount, &nu.LimitedUse, &nu.RequiresCode, &nu.DiscountType, &nu.Code, &nu.ExpiresOn, &nu.StartsOn, &nu.NumberOfUses, &nu.Name, &nu.LoginRequired).Scan(&createdID, &createdAt)
	return createdID, createdAt, err
}

const discountUpdateQuery = `
    UPDATE discounts
    SET
        amount = $1, 
        limited_use = $2, 
        requires_code = $3, 
        discount_type = $4, 
        code = $5, 
        expires_on = $6, 
        updated_on = NOW()
        starts_on = $8, 
        number_of_uses = $9, 
        name = $10, 
        login_required = $11
    WHERE id = $11
    RETURNING updated_on;
`

func (pg *postgres) UpdateDiscount(db storage.Querier, updated *models.Discount) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(discountUpdateQuery, &updated.Amount, &updated.LimitedUse, &updated.RequiresCode, &updated.DiscountType, &updated.Code, &updated.ExpiresOn, &updated.StartsOn, &updated.NumberOfUses, &updated.Name, &updated.LoginRequired, &updated.ID).Scan(&t)
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
