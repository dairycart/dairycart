package pg_test

import (
	"fmt"

	"github.com/go-pg/pg"
)

type Params struct {
	X int
	Y int
}

func (p *Params) Sum() int {
	return p.X + p.Y
}

// go-pg recognizes placeholders (`?`) in queries and replaces them
// with parameters when queries are executed. Parameters are escaped
// before replacing according to PostgreSQL rules. Specifically:
//   - all parameters are properly quoted against SQL injections;
//   - null byte is removed;
//   - JSON/JSONB gets `\u0000` escaped as `\\u0000`.
func Example_placeholders() {
	var num int

	// Simple params.
	_, err := db.Query(pg.Scan(&num), "SELECT ?", 42)
	if err != nil {
		panic(err)
	}
	fmt.Println("simple:", num)

	// Indexed params.
	_, err = db.Query(pg.Scan(&num), "SELECT ?0 + ?0", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println("indexed:", num)

	// Named params.
	params := &Params{
		X: 1,
		Y: 1,
	}
	_, err = db.Query(pg.Scan(&num), "SELECT ?x + ?y + ?Sum", params)
	if err != nil {
		panic(err)
	}
	fmt.Println("named:", num)

	// Global params.
	_, err = db.WithParam("z", 1).Query(pg.Scan(&num), "SELECT ?x + ?y + ?z", params)
	if err != nil {
		panic(err)
	}
	fmt.Println("global:", num)

	// Output: simple: 42
	// indexed: 2
	// named: 4
	// global: 3
}
