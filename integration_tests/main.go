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

func buildPath(parts ...string) string {
	return fmt.Sprintf("%s/%s", baseURL, strings.Join(parts, "/"))
}

func mapToQueryValues(in map[string]string) string {
	out := url.Values{}
	for k, v := range in {
		out.Set(k, v)
	}
	return out.Encode()
}

func buildURL(path string, queryParams map[string]string) string {
	url, _ := url.Parse(path)
	queryString := mapToQueryValues(queryParams)
	url.RawQuery = queryString
	return url.String()
}

func ensureThatDairycartIsAlive() error {
	path := buildPath("health")
	url := buildURL(path, nil)
	dairyCartIsDown := true
	numberOfAttempts := 0
	for dairyCartIsDown {
		_, err := http.Get(url)
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
	path := buildPath("product", sku)
	url := buildURL(path, nil)
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func retrieveProduct(sku string) (*http.Response, error) {
	path := buildPath("product", sku)
	url := buildURL(path, nil)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func retrieveListOfProducts(queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("products")
	url := buildURL(path, queryFilter)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func createProduct(JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	path := buildPath("product")
	url := buildURL(path, nil)
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func updateProduct(sku string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	path := buildPath("product", sku)
	url := buildURL(path, nil)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}

func deleteProduct(sku string) (*http.Response, error) {
	path := buildPath("product", sku)
	url := buildURL(path, nil)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}
