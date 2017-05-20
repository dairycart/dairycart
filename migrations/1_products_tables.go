package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-pg/migrations"
)

func init() {
	migrationName := "1_products.sql"

	migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating products tables...")
		fileName := fmt.Sprintf("sql/up/%s", migrationName)
		migration, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(string(migration))
		return err
	}, func(db migrations.DB) error {
		fileName := fmt.Sprintf("sql/down/%s", migrationName)
		migration, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}

		_, err = db.Exec(string(migration))
		return err
	})
}
