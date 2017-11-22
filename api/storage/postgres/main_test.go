package postgres

import (
	"testing"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/Masterminds/squirrel"
	"github.com/stretchr/testify/require"
)

func TestApplyQueryFilterToQueryBuilder(t *testing.T) {
	t.Parallel()
	baseQueryBuilder := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("things").
		From("stuff").
		Where(squirrel.Eq{"condition": true})

	t.Run("basic usecase", func(*testing.T) {
		exampleQF := &models.QueryFilter{
			Limit: 15,
			Page:  2,
		}
		expected := `SELECT things FROM stuff WHERE condition = $1 LIMIT 15`

		x := applyQueryFilterToQueryBuilder(baseQueryBuilder, exampleQF, false)
		actual, args, err := x.ToSql()
		require.Equal(t, expected, actual, "expected and actual queries don't match")
		require.Nil(t, err)
	 	require.NotEmpty(t, args)
	})

	t.Run("returns query builder if query filter is nil", func(*testing.T) {
		expected := `SELECT things FROM stuff WHERE condition = $1`

		x := applyQueryFilterToQueryBuilder(baseQueryBuilder, nil, false)
		actual, args, err := x.ToSql()
		require.Equal(t, expected, actual, "expected and actual queries don't match")
		require.Nil(t, err)
		require.NotEmpty(t, args)
	})

	t.Run("whole kit and kaboodle", func(*testing.T){
		exampleQF := &models.QueryFilter{
			Limit: 20,
			Page:  6,
			CreatedAfter: time.Now(),
			CreatedBefore: time.Now(),
			UpdatedAfter: time.Now(),
			UpdatedBefore: time.Now(),
		}
		expected := `SELECT things FROM stuff WHERE condition = $1 AND created_on > $2 AND created_on < $3 AND updated_on > $4 AND updated_on < $5 LIMIT 20 OFFSET 100`

		x := applyQueryFilterToQueryBuilder(baseQueryBuilder, exampleQF, true)
		actual, args, err := x.ToSql()
		require.Equal(t, expected, actual, "expected and actual queries don't match")
		require.Nil(t, err)
		require.NotEmpty(t, args)
	})

	t.Run("with zero limit", func(*testing.T) {
		exampleQF := &models.QueryFilter{
			Limit: 0,
			Page:  1,
		}
		expected := `SELECT things FROM stuff WHERE condition = $1 LIMIT 25`

		x := applyQueryFilterToQueryBuilder(baseQueryBuilder, exampleQF, false)
		actual, args, err := x.ToSql()
		require.Equal(t, expected, actual, "expected and actual queries don't match")
		require.Nil(t, err)
		require.NotEmpty(t, args)
	})

}
