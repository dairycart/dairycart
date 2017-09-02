package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type Order struct {
	ID        uint64    `json:"id"`
	CreatedOn time.Time `json:"created_on"`
}

type OrdersPage struct {
	Page
	Orders []Order
}

type OrderPage struct {
	Page
	Order Order
}

func serveOrder(w http.ResponseWriter, r *http.Request) {
	p := &OrderPage{
		Page: Page{Title: "Products"},
		Order: Order{
			ID:        1,
			CreatedOn: time.Now(),
		},
	}

	baseTemplatePath := filepath.Join(templateDir, "base.html")
	innerTemplatePath := filepath.Join(templateDir, "order.html")

	tmpl, err := template.ParseFiles(baseTemplatePath, innerTemplatePath)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", p); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func serveOrders(w http.ResponseWriter, r *http.Request) {
	p := &OrdersPage{
		Page: Page{Title: "Products"},
		Orders: []Order{
			{
				ID:        1,
				CreatedOn: time.Now(),
			},
			{
				ID:        2,
				CreatedOn: time.Now(),
			},
			{
				ID:        3,
				CreatedOn: time.Now(),
			},
			{
				ID:        4,
				CreatedOn: time.Now(),
			},
			{
				ID:        5,
				CreatedOn: time.Now(),
			},
		},
	}

	baseTemplatePath := filepath.Join(templateDir, "base.html")
	innerTemplatePath := filepath.Join(templateDir, "orders.html")

	tmpl, err := template.ParseFiles(baseTemplatePath, innerTemplatePath)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "base", p); err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
