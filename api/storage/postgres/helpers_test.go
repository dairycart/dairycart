package postgres

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func formatQueryForSQLMock(query string) string {
	for _, x := range []string{"$", "(", ")", "=", "*", ".", "+", "?", ",", "-"} {
		query = strings.Replace(query, x, fmt.Sprintf(`\%s`, x), -1)
	}
	return query
}

func generateExampleTimeForTests(t *testing.T) time.Time {
	t.Helper()
	out, err := time.Parse("2006-01-02 03:04:00.000000", "2016-12-31 12:00:00.000000")
	require.Nil(t, err)
	return out
}
