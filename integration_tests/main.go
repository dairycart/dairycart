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

var (
	baseURL   string
	requester *Requester
)

const (
	maxAttempts       = 10
	currentAPIVersion = `v1`
	validPassword     = "Pa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rdPa$$w0rd"
)

// Requester is what we use to make requests
type Requester struct {
	http.Client
	AdminAuthCookie   *http.Cookie
	RegularAuthCookie *http.Cookie
}

// ExecuteRequestAsAdmin takes a regular prepared Request struct and adds our admin Auth cookie before making the request.
func (r *Requester) ExecuteRequestAsAdmin(req *http.Request) (*http.Response, error) {
	req.AddCookie(r.AdminAuthCookie)
	return r.Do(req)
}

// ExecuteRequestAsRegularUser takes a regular prepared Request struct and adds our regular Auth cookie before making the request.
func (r *Requester) ExecuteRequestAsRegularUser(req *http.Request) (*http.Response, error) {
	req.AddCookie(r.RegularAuthCookie)
	return r.Do(req)
}

func init() {
	baseURL = os.Getenv("INTEGRATION_API_URL")
	dbURL := os.Getenv("DAIRYCART_DB_URL")
	if dbURL == "" {
		// running outside of docker, and therefore debugging
		baseURL = `http://localhost`
		dbURL = "postgres://dairycart:hunter2@localhost:2345/dairycart?sslmode=disable"
	}

	ensureThatDairycartIsAlive()
	createUsers(dbURL)
	requester = &Requester{}
	getUserCookies()
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
	u, _ := url.Parse(path)
	queryString := mapToQueryValues(queryParams)
	u.RawQuery = queryString
	return u.String()
}

func createUsers(dbURL string) {
	// Connect to the database
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

	_, err = db.Exec(`INSERT INTO users ("first_name", "last_name", "username", "email", "password", "salt", "is_admin") VALUES ('admin', 'user', 'regular', 'regular@user.com', '$2a$13$NSHE6gf1FlATM3YUgVWGIe9Ao4DUUHydreuE7eZoc8DNbxq1rw.yq', 'fake_salt_here'::bytea, 'false')`)
	if err != nil {
		log.Fatalf("error encountered creating regular user: %v", err)
	}
}

func getUserCookies() {
	resp, err := loginUser("regular", validPassword)
	if err != nil {
		log.Fatal(err)
	}
	requester.RegularAuthCookie = resp.Cookies()[0]

	resp, err = loginUser("admin", validPassword)
	if err != nil {
		log.Fatal(err)
	}
	requester.AdminAuthCookie = resp.Cookies()[0]
}

func ensureThatDairycartIsAlive() {
	path := buildPath("health")
	u := buildURL(path, nil)
	dairyCartIsDown := true
	numberOfAttempts := 0
	for dairyCartIsDown {
		_, err := http.Get(u)
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

func createNewUser(JSONBody string, createAsSuperUser bool) (*http.Response, error) {
	u := buildVersionlessPath("user")
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, u, body)
	if createAsSuperUser {
		return requester.ExecuteRequestAsAdmin(req)
	}
	return requester.ExecuteRequestAsRegularUser(req)
}

func deleteUser(userID string, deleteAsSuperUser bool) (*http.Response, error) {
	u := buildPath("user", userID)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	if deleteAsSuperUser {
		return requester.ExecuteRequestAsAdmin(req)
	}
	return requester.ExecuteRequestAsRegularUser(req)
}

func loginUser(username string, password string) (*http.Response, error) {
	u := buildVersionlessPath("login")
	body := strings.NewReader(fmt.Sprintf(`
		{
			"username": "%s",
			"password": "%s"
		}
	`, username, password))
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return requester.Do(req)
}

func logoutUser(cookie *http.Cookie) (*http.Response, error) {
	u := buildVersionlessPath("logout")
	req, _ := http.NewRequest(http.MethodPost, u, nil)
	req.AddCookie(cookie)
	return requester.Do(req)
}

////////////////////////////////////////////////////////
//                                                    //
//                 Product Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func checkProductExistence(sku string) (*http.Response, error) {
	u := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodHead, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

func retrieveProduct(sku string) (*http.Response, error) {
	u := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

func retrieveListOfProducts(queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("products")
	u := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

func createProduct(JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u := buildPath("product")
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return requester.ExecuteRequestAsRegularUser(req)
}

func updateProduct(sku string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return requester.ExecuteRequestAsRegularUser(req)
}

func deleteProduct(sku string) (*http.Response, error) {
	u := buildPath("product", sku)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

////////////////////////////////////////////////////////
//                                                    //
//              Product Root Functions                //
//                                                    //
////////////////////////////////////////////////////////

func retrieveProductRoot(rootID string) (*http.Response, error) {
	u := buildPath("product_root", rootID)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

func retrieveProductRoots(queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("product_roots")
	u := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

func deleteProductRoot(rootID string) (*http.Response, error) {
	u := buildPath("product_root", rootID)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

////////////////////////////////////////////////////////
//                                                    //
//             Product Option Functions               //
//                                                    //
////////////////////////////////////////////////////////

func retrieveProductOptions(productID string, queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("product", productID, "options")
	u := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

func createProductOptionForProduct(productID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u := buildPath("product", productID, "options")
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return requester.ExecuteRequestAsRegularUser(req)
}

func updateProductOption(optionID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u := buildPath("product_options", optionID)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return requester.ExecuteRequestAsRegularUser(req)
}

func deleteProductOption(optionID string) (*http.Response, error) {
	u := buildPath("product_options", optionID)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

////////////////////////////////////////////////////////
//                                                    //
//          Product Option Value Functions            //
//                                                    //
////////////////////////////////////////////////////////

func createProductOptionValueForOption(optionID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u := buildPath("product_options", optionID, "value")
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return requester.ExecuteRequestAsRegularUser(req)
}

func updateProductOptionValueForOption(valueID string, JSONBody string) (*http.Response, error) {
	body := strings.NewReader(JSONBody)
	u := buildPath("product_option_values", valueID)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return requester.ExecuteRequestAsRegularUser(req)
}

func deleteProductOptionValueForOption(optionID string) (*http.Response, error) {
	u := buildPath("product_option_values", optionID)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

////////////////////////////////////////////////////////
//                                                    //
//                Discount Functions                  //
//                                                    //
////////////////////////////////////////////////////////

func getDiscountByID(discountID string) (*http.Response, error) {
	u := buildPath("discount", discountID)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

func getListOfDiscounts(queryFilter map[string]string) (*http.Response, error) {
	path := buildPath("discounts")
	u := buildURL(path, queryFilter)
	req, _ := http.NewRequest(http.MethodGet, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}

func createDiscount(JSONBody string) (*http.Response, error) {
	u := buildPath("discount")
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPost, u, body)
	return requester.ExecuteRequestAsRegularUser(req)
}

func updateDiscount(discountID string, JSONBody string) (*http.Response, error) {
	u := buildPath("discount", discountID)
	body := strings.NewReader(JSONBody)
	req, _ := http.NewRequest(http.MethodPatch, u, body)
	return requester.ExecuteRequestAsRegularUser(req)
}

func deleteDiscount(discountID string) (*http.Response, error) {
	u := buildPath("discount", discountID)
	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	return requester.ExecuteRequestAsRegularUser(req)
}
