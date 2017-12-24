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
        name,
        discount_type,
        amount,
        expires_on,
        requires_code,
        code,
        limited_use,
        number_of_uses,
        login_required,
        starts_on,
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
	err := db.QueryRow(discountQueryByCode, code).Scan(&d.ID, &d.Name, &d.DiscountType, &d.Amount, &d.ExpiresOn, &d.RequiresCode, &d.Code, &d.LimitedUse, &d.NumberOfUses, &d.LoginRequired, &d.StartsOn, &d.CreatedOn, &d.UpdatedOn, &d.ArchivedOn)
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
        expires_on,
        requires_code,
        code,
        limited_use,
        number_of_uses,
        login_required,
        starts_on,
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

	err := db.QueryRow(discountSelectionQuery, id).Scan(&d.ID, &d.Name, &d.DiscountType, &d.Amount, &d.ExpiresOn, &d.RequiresCode, &d.Code, &d.LimitedUse, &d.NumberOfUses, &d.LoginRequired, &d.StartsOn, &d.CreatedOn, &d.UpdatedOn, &d.ArchivedOn)

	return d, err
}

func buildDiscountListRetrievalQuery(qf *models.QueryFilter) (string, []interface{}) {
	sqlBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	queryBuilder := sqlBuilder.
		Select(
			"id",
			"name",
			"discount_type",
			"amount",
			"expires_on",
			"requires_code",
			"code",
			"limited_use",
			"number_of_uses",
			"login_required",
			"starts_on",
			"created_on",
			"updated_on",
			"archived_on",
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
			&d.Name,
			&d.DiscountType,
			&d.Amount,
			&d.ExpiresOn,
			&d.RequiresCode,
			&d.Code,
			&d.LimitedUse,
			&d.NumberOfUses,
			&d.LoginRequired,
			&d.StartsOn,
			&d.CreatedOn,
			&d.UpdatedOn,
			&d.ArchivedOn,
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
            name, discount_type, amount, expires_on, requires_code, code, limited_use, number_of_uses, login_required, starts_on
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

	err := db.QueryRow(discountCreationQuery, &nu.Name, &nu.DiscountType, &nu.Amount, &nu.ExpiresOn, &nu.RequiresCode, &nu.Code, &nu.LimitedUse, &nu.NumberOfUses, &nu.LoginRequired, &nu.StartsOn).Scan(&createdID, &createdAt)
	return createdID, createdAt, err
}

const discountUpdateQuery = `
    UPDATE discounts
    SET
        name = $1,
        discount_type = $2,
        amount = $3,
        expires_on = $4,
        requires_code = $5,
        code = $6,
        limited_use = $7,
        number_of_uses = $8,
        login_required = $9,
        starts_on = $10,
        updated_on = NOW()
    WHERE id = $11
    RETURNING updated_on;
`

func (pg *postgres) UpdateDiscount(db storage.Querier, updated *models.Discount) (time.Time, error) {
	var t time.Time
	err := db.QueryRow(discountUpdateQuery, &updated.Name, &updated.DiscountType, &updated.Amount, &updated.ExpiresOn, &updated.RequiresCode, &updated.Code, &updated.LimitedUse, &updated.NumberOfUses, &updated.LoginRequired, &updated.StartsOn, &updated.ID).Scan(&t)
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
