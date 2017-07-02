package dairytest

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var requester *Requester

const (
	maxAttempts = 10
	baseURL     = `http://dairycart/v1`
	password    = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
)

// Requester is what we use to make requests
type Requester struct {
	http.Client
	AuthToken string
}

func (r *Requester) execRequest(req *http.Request) (*http.Response, error) {
	// authorization step goes here
	return r.Do(req)
}

func init() {
	ensureThatDairycartIsAlive()
	requester = &Requester{}
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

func ensureThatDairycartIsAlive() {
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
				log.Fatalf("Maximum number of attempts made, something's gone awry")
			}
		} else {
			dairyCartIsDown = false
		}
	}
}

func createNewUser(JSONBody string) (*http.Response, error) {
	url := `http://dairycart/user`
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.Do(req)
}

func authenticateUser(email string, password string) (*http.Response, error) {
	url := `http://dairycart/login`
	body := strings.NewReader(fmt.Sprintf(`
		{
			"email": "%s",
			"password": "%s"
		}
	`, email, password))
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.Do(req)
}

func checkProductExistence(sku string) (*http.Response, error) {
	url := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodHead, url, nil)
	return requester.execRequest(req)
}

func retrieveProduct(sku string) (*http.Response, error) {
	url := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.execRequest(req)
}

func retrieveListOfProducts(queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("products")
	url := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.execRequest(req)
}

func createProduct(JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product")
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.execRequest(req)
}

func updateProduct(sku string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodPut, url, body)
	return requester.execRequest(req)
}

func deleteProduct(sku string) (*http.Response, error) {
	url := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	return requester.execRequest(req)
}

func retrieveProductOptions(productID string, queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("product", productID, "options")
	url := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.execRequest(req)
}

func createProductOptionForProduct(productID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product", productID, "options")
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.execRequest(req)
}

func updateProductOption(optionID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product_options", optionID)
	req, _ := http.NewRequest(http.MethodPut, url, body)
	return requester.execRequest(req)
}

func createProductOptionValueForOption(optionID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product_options", optionID, "value")
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.execRequest(req)
}

func updateProductOptionValueForOption(valueID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product_option_values", valueID)
	req, _ := http.NewRequest(http.MethodPut, url, body)
	return requester.execRequest(req)
}

func getDiscountByID(discountID string) (*http.Response, error) {
	url := buildPath("discount", discountID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.execRequest(req)
}

func getListOfDiscounts(queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("discounts")
	url := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.execRequest(req)
}

func createDiscount(JSONBody string) (*http.Response, error) {
	url := buildPath("discount")
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.execRequest(req)
}

func updateDiscount(discountID string, JSONBody string) (*http.Response, error) {
	url := buildPath("discount", discountID)
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPut, url, body)
	return requester.execRequest(req)
}
