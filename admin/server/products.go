package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func subset(ps []Product, start, end int) []Product {
	log.Printf("subset called with start: %d, end: %d\n", start, end)
	return ps[start:end]
}

// Product describes something a user can buy
type Product struct {
	ID          uint64    `json:"id"`
	CreatedOn   time.Time `json:"created_on"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	SKU         string    `json:"sku"`
	Quantity    uint32    `json:"quantity"`
	Price       float32   `json:"price"`
}

type ProductsPage struct {
	Page
	ProductChunks [][]Product
}

func serveProducts(w http.ResponseWriter, r *http.Request) {
	p := &ProductsPage{
		Page: Page{Title: "Products"},
		ProductChunks: [][]Product{
			{
				{

					ID:        1,
					CreatedOn: time.Now(),
					Name:      "Farts",
					Price:     123.45,
					Quantity:  321,
				},
				{
					ID:        2,
					CreatedOn: time.Now(),
					Name:      "Butts",
					Price:     43.21,
					Quantity:  420,
				},
				{
					ID:        3,
					CreatedOn: time.Now(),
					Name:      "Dongs",
					Price:     12.34,
					Quantity:  666,
				},
				{
					ID:        4,
					CreatedOn: time.Now(),
					Name:      "Shit",
					Price:     12.34,
					Quantity:  666,
				},
			},
			{
				{

					ID:        1,
					CreatedOn: time.Now(),
					Name:      "Farts",
					Price:     123.45,
					Quantity:  321,
				},
				{
					ID:        2,
					CreatedOn: time.Now(),
					Name:      "Butts",
					Price:     43.21,
					Quantity:  420,
				},
				{
					ID:        3,
					CreatedOn: time.Now(),
					Name:      "Dongs",
					Price:     12.34,
					Quantity:  666,
				},
			},
		},
	}

	baseTemplatePath := filepath.Join(templateDir, "base.html")
	innerTemplatePath := filepath.Join(templateDir, "products.html")

	fm := template.FuncMap{
		"subset": subset,
	}

	tmpl, err := template.ParseFiles(baseTemplatePath, innerTemplatePath)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	tmpl = tmpl.Funcs(fm)

	if err := tmpl.ExecuteTemplate(w, "base", p); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
