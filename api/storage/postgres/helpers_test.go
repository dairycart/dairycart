package postgres

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func formatQueryForSQLMock(query string) string {
	for _, x := range []string{"$", "(", ")", "=", "*", ".", "+", "?", ",", "-"} {
		query = strings.Replace(query, x, fmt.Sprintf(`\%s`, x), -1)
	}
	return query
}

func generateExampleTimeForTests() time.Time {
	out, err := time.Parse("2006-01-02 03:04:00.000000", "2016-12-31 12:00:00.000000")
	if err != nil {
		log.Fatalf("error parsing time")
	}
	return out
}
