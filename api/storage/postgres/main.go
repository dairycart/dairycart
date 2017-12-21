package postgres

import (
	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/Masterminds/squirrel"
)

var _ storage.Storer = (*postgres)(nil)

type postgres struct{}

func NewPostgres() *postgres {
	return &postgres{}
}

func applyQueryFilterToQueryBuilder(queryBuilder squirrel.SelectBuilder, qf *models.QueryFilter, includeOffset bool) squirrel.SelectBuilder {
	if qf == nil {
		return queryBuilder
	}

	if qf.Limit > 0 {
		queryBuilder = queryBuilder.Limit(uint64(qf.Limit))
	} else {
		queryBuilder = queryBuilder.Limit(25)
	}

	if qf.Page > 1 && includeOffset {
		offset := (qf.Page - 1) * uint64(qf.Limit)
		queryBuilder = queryBuilder.Offset(offset)
	}

	if !qf.CreatedAfter.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Gt{"created_on": qf.CreatedAfter})
	}

	if !qf.CreatedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"created_on": qf.CreatedBefore})
	}

	if !qf.UpdatedAfter.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Gt{"updated_on": qf.UpdatedAfter})
	}

	if !qf.UpdatedBefore.IsZero() {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"updated_on": qf.UpdatedBefore})
	}

	if !qf.IncludeArchived {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"archived_on": nil})
	}

	return queryBuilder
}
