package dairytest

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	// PQ is the postgres driver
	_ "github.com/lib/pq"
)

var requester *Requester

const (
	maxAttempts       = 10
	baseURL           = `http://dairycart`
	currentAPIVersion = `v1`
	validPassword     = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
)

// Requester is what we use to make requests
type Requester struct {
	http.Client
	AuthCookie *http.Cookie
}

// ExecuteAuthorizedRequest takes a regular prepared Request struct and adds our superuser Auth cookie before making the request.
func (r *Requester) ExecuteAuthorizedRequest(req *http.Request) (*http.Response, error) {
	req.AddCookie(r.AuthCookie)
	return r.Do(req)
}

func init() {
	ensureThatDairycartIsAlive()
	createSuperUser()
	requester = &Requester{}
	getSuperUserCookie()
}

////////////////////////////////////////////////////////
//                                                    //
//                 Helper Functions                   //
//                                                    //
////////////////////////////////////////////////////////

func buildPath(parts ...string) string {
	return fmt.Sprintf("%s/%s/%s", baseURL, currentAPIVersion, strings.Join(parts, "/"))
}

func buildVersionlessPath(parts ...string) string {
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

func createSuperUser() {
	// Connect to the database
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("error encountered connecting to database: %v", err)
	}

	// this is here because sometimes I run Dairycart locally and debug individual tests, and the below insert statement will fail without this step.
	_, err = db.Exec(`DELETE FROM users WHERE id IS NOT NULL`)
	if err != nil {
		log.Fatalf("error encountered deleting existing users: %v", err)
	}

	_, err = db.Exec(`INSERT INTO users ("first_name", "last_name", "username", "email", "password", "salt", "is_admin") VALUES ('admin', 'user', 'admin', 'admin@user.com', '$2a$13$NSHE6gf1FlATM3YUgVWGIe9Ao4DUUHydreuE7eZoc8DNbxq1rw.yq', 'fake_salt_here'::bytea, 'true')`)
	if err != nil {
		log.Fatalf("error encountered creating super user: %v", err)
	}
}

func getSuperUserCookie() {
	resp, err := loginUser("admin", validPassword)
	if err != nil {
		log.Fatal(err)
	}
	requester.AuthCookie = resp.Cookies()[0]
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

////////////////////////////////////////////////////////
//                                                    //
//                  Auth Functions                    //
//                                                    //
////////////////////////////////////////////////////////

func createNewUser(JSONBody string) (*http.Response, error) {
	url := buildVersionlessPath("user")
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.ExecuteAuthorizedRequest(req)
}

func loginUser(username string, password string) (*http.Response, error) {
	url := buildVersionlessPath("login")
	body := strings.NewReader(fmt.Sprintf(`
		{
			"username": "%s",
			"password": "%s"
		}
	`, username, password))
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.Do(req)
}

func logoutUser(username string, password string) (*http.Response, error) {
	url := buildVersionlessPath("logout")
	req, _ := http.NewRequest(http.MethodPost, url, nil)
	return requester.Do(req)
}

func deleteUser(userID string) (*http.Response, error) {
	url := buildPath("user", userID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	return requester.Do(req)
}

////////////////////////////////////////////////////////
//                                                    //
//                 Product Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func checkProductExistence(sku string) (*http.Response, error) {
	url := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodHead, url, nil)
	return requester.Do(req)
}

func retrieveProduct(sku string) (*http.Response, error) {
	url := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.Do(req)
}

func retrieveListOfProducts(queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("products")
	url := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.Do(req)
}

func createProduct(JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product")
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.Do(req)
}

func updateProduct(sku string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodPatch, url, body)
	return requester.Do(req)
}

func deleteProduct(sku string) (*http.Response, error) {
	url := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	return requester.Do(req)
}

////////////////////////////////////////////////////////
//                                                    //
//             Product Option Functions               //
//                                                    //
////////////////////////////////////////////////////////

func retrieveProductOptions(productID string, queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("product", productID, "options")
	url := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.Do(req)
}

func createProductOptionForProduct(productID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product", productID, "options")
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.Do(req)
}

func updateProductOption(optionID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product_options", optionID)
	req, _ := http.NewRequest(http.MethodPatch, url, body)
	return requester.Do(req)
}

func deleteProductOption(optionID string) (*http.Response, error) {
	url := buildPath("product_options", optionID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	return requester.Do(req)
}

////////////////////////////////////////////////////////
//                                                    //
//          Product Option Value Functions            //
//                                                    //
////////////////////////////////////////////////////////

func createProductOptionValueForOption(optionID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product_options", optionID, "value")
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.Do(req)
}

func updateProductOptionValueForOption(valueID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	url := buildPath("product_option_values", valueID)
	req, _ := http.NewRequest(http.MethodPatch, url, body)
	return requester.Do(req)
}

func deleteProductOptionValueForOption(optionID string) (*http.Response, error) {
	url := buildPath("product_option_values", optionID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	return requester.Do(req)
}

////////////////////////////////////////////////////////
//                                                    //
//                Discount Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func getDiscountByID(discountID string) (*http.Response, error) {
	url := buildPath("discount", discountID)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.Do(req)
}

func getListOfDiscounts(queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("discounts")
	url := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	return requester.Do(req)
}

func createDiscount(JSONBody string) (*http.Response, error) {
	url := buildPath("discount")
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, url, body)
	return requester.Do(req)
}

func updateDiscount(discountID string, JSONBody string) (*http.Response, error) {
	url := buildPath("discount", discountID)
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPatch, url, body)
	return requester.Do(req)
}

func deleteDiscount(discountID string) (*http.Response, error) {
	url := buildPath("discount", discountID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	return requester.Do(req)
}
