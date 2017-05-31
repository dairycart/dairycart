package dairytest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var client *http.Client

const (
	maxAttempts = 10
	baseURL     = `http://dairycart/v1`
)

func init() {
	client = &http.Client{}
}

func buildURL(parts ...string) *url.URL {
	url, _ := url.Parse(fmt.Sprintf("%s/%s", baseURL, strings.Join(parts, "/")))
	return url
}

func ensureThatDairycartIsAlive() error {
	url := buildURL("health")
	dairyCartIsDown := true
	numberOfAttempts := 0
	for dairyCartIsDown {
		_, err := http.Get(url.String())
		if err != nil {
			log.Printf("waiting half a second before pinging Dairycart again")
			time.Sleep(500 * time.Millisecond)
			numberOfAttempts++
			if numberOfAttempts >= maxAttempts {
				return errors.New("Maximum number of attempts made, something's gone awry")
			}
		} else {
			dairyCartIsDown = false
		}
	}
	return nil
}

func checkProductExistence(sku string) (*http.Response, error) {
	url := buildURL("product", sku)
	req, err := http.NewRequest(http.MethodHead, url.String(), nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func retrieveProduct(sku string) (*http.Response, error) {
	url := buildURL("product", sku)
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func retrieveListOfProducts() (*http.Response, error) {
	url := buildURL("products")
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func createProduct(JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildURL("product")
	req, err := http.NewRequest(http.MethodPost, url.String(), body)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func updateProduct(sku string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildURL("product", sku)
	req, err := http.NewRequest(http.MethodPut, url.String(), body)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func deleteProduct(sku string) (*http.Response, error) {
	url := buildURL("product", sku)
	req, err := http.NewRequest(http.MethodDelete, url.String(), nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}
