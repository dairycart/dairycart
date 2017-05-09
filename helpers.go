package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-pg/pg/orm"
)

// AddQueryParams takes some URL values and the result of a db.Model() call
// and adds the appropriate conditions
func AddQueryParams(v url.Values, q *orm.Query) {
	for k, v := range v {
		q.Where(fmt.Sprintf("%s = ?", k), v[0])
	}
}

// GenericListHandler generalizes list handling
func GenericListHandler(listOfThings *interface{}, res http.ResponseWriter, req *http.Request) {
	pluralModel := db.Model(listOfThings)
	AddQueryParams(req.URL.Query(), pluralModel)
	err := pluralModel.Select()
	if err != nil {
		log.Printf("Error encountered querying for products: %v", err)
	}
	json.NewEncoder(res).Encode(pluralModel)
}
