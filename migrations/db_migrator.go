package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

func main() {
	flag.Parse()
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	dbOptions, err := pg.ParseURL(dbURL)
	if err != nil {
		log.Fatalf("While parsing database URL: %v, encountered error: %v", dbURL, err)
	}
	db := pg.Connect(dbOptions)

	oldVersion, newVersion, err := migrations.Run(db, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		os.Exit(1)
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	}
}
