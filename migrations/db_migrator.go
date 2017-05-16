package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
)

const verbose = true

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
		exitf(err.Error())
	}
	if verbose {
		if newVersion != oldVersion {
			fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
		}
	}
}

func errorf(s string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, s+"\n", args...)
}

func exitf(s string, args ...interface{}) {
	errorf(s, args...)
	os.Exit(1)
}
