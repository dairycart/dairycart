package dairytest

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var client *http.Client

const (
	maxAttempts = 10
	baseURL     = `http://dairycart`
)

func init() {
	client = &http.Client{}
}

func buildURL(parts ...string) string {
	return fmt.Sprintf("%s/%s", baseURL, strings.Join(parts, "/"))
}

func ensureThatDairycartIsAlive() error {
	url := buildURL("health")
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
	url := buildURL("product", sku)
	return http.Head(url)
}

func retrieveProduct(sku string) (*http.Response, error) {
	url := buildURL("product", sku)
	return http.Get(url)
}

func retrieveListOfProducts() (*http.Response, error) {
	url := buildURL("products")
	return http.Get(url)
}

func updateProduct(sku string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildURL("product", sku)
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}

	return client.Do(req)
}

func deleteProduct(sku string) (*http.Response, error) {
	url := buildURL("product", sku)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		log.Fatalf("failed to build request: %v", err)
	}
	return client.Do(req)
}
