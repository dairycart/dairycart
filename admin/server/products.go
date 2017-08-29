package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi"

	"github.com/dairycart/dairyclient/v1"
)

const (
	maxGallerySize = 5
)

type ProductsPage struct {
	Page
	ProductChunks [][]dairyclient.Product
}

type ProductPage struct {
	Page
	Product *dairyclient.Product
}

func serveProduct(res http.ResponseWriter, req *http.Request) {
	sku := chi.URLParam(req, "sku")
	dairyClient, err := buildClientFromRequest(res, req)
	if err != nil {
		return
	}

	product, err := dairyClient.GetProduct(sku)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	p := &ProductPage{
		Page:    Page{Title: "Products"},
		Product: product,
	}

	baseTemplatePath := filepath.Join(templateDir, "base.html")
	innerTemplatePath := filepath.Join(templateDir, "product.html")

	tmpl, err := template.ParseFiles(baseTemplatePath, innerTemplatePath)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(res, "base", p); err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func splitProducts(in []dairyclient.Product) [][]dairyclient.Product {
	var out [][]dairyclient.Product

	x := []dairyclient.Product{}
	for i, p := range in {
		if i%maxGallerySize == 0 {
			out = append(out, x)
			x = []dairyclient.Product{}
		}
		x = append(x, p)
	}
	out = append(out, x)
	return out
}

func serveProducts(res http.ResponseWriter, req *http.Request) {
	dairyClient, err := buildClientFromRequest(res, req)
	if err != nil {
		// TODO: this pattern is bad and you should feel bad
		return
	}

	products, err := dairyClient.GetProducts(nil)
	if err != nil {
		log.Printf("error encountered retrieving products with the dairyclient: %v\n", err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	productChunks := splitProducts(products)

	p := &ProductsPage{
		Page:          Page{Title: "Products"},
		ProductChunks: productChunks,
	}

	baseTemplatePath := filepath.Join(templateDir, "base.html")
	innerTemplatePath := filepath.Join(templateDir, "products.html")

	tmpl, err := template.ParseFiles(baseTemplatePath, innerTemplatePath)
	if err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(res, "base", p); err != nil {
		log.Println(err)
		http.Error(res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
