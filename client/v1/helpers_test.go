package dairyclient_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/dairycart/dairycart/client/v1"
	"github.com/dairycart/dairycart/models/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tdewolff/minify"
	jsonMinify "github.com/tdewolff/minify/json"
)

const (
	exampleID       = 666
	exampleURL      = `http://www.dairycart.com`
	exampleUsername = `username`
	examplePassword = `password` // lol not really
	exampleSKU      = `sku`
	exampleBadJSON  = `{"invalid lol}`
	timeLayout      = "2006-01-02T15:04:05.000000Z"
)

////////////////////////////////////////////////////////
//                                                    //
//               Test Helper Functions                //
//                                                    //
////////////////////////////////////////////////////////

func buildTestTime(t *testing.T) time.Time {
	t.Helper()
	xt, err := time.Parse(timeLayout, "2017-12-10T15:58:43.136458Z")
	require.NoError(t, err)
	return xt
}

func buildTestDairytime(t *testing.T) *models.Dairytime {
	t.Helper()
	xt, err := time.Parse(timeLayout, "2017-12-10T15:58:43.136458Z")
	require.NoError(t, err)
	return &models.Dairytime{Time: xt}
}

func buildTestCookie() *http.Cookie {
	c := &http.Cookie{Name: "dairycart"}
	return c
}

func buildTestClient(t *testing.T, ts *httptest.Server) *dairyclient.V1Client {
	t.Helper()

	u, err := url.Parse(ts.URL)
	assert.Nil(t, err)

	c := &dairyclient.V1Client{
		URL:        u,
		Client:     ts.Client(),
		AuthCookie: buildTestCookie(),
	}

	return c
}

func loadExampleResponse(t *testing.T, name string) string {
	t.Helper()
	data, err := ioutil.ReadFile(fmt.Sprintf("example_responses/%s.json", name))
	if err != nil {
		log.Printf("error encountered reading example response file: %v\n", err)
		t.FailNow()
	}
	return string(data)
}

func obligatoryLoginHandler(addCookie bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if addCookie {
			cookie := &http.Cookie{
				Name: "dairycart",
			}
			http.SetCookie(w, cookie)
		}
	})
}

func handlerGenerator(handlers map[string]http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if handlerFunc, ok := handlers[r.URL.Path]; ok {
			handlerFunc(w, r)
			return
		} else {
			http.NotFound(w, r)
			return
		}
	})
}

func minifyJSON(t *testing.T, jsonBody string) string {
	t.Helper()

	jsonMinifier := minify.New()
	jsonMinifier.AddFunc("application/json", jsonMinify.Minify)
	minified, err := jsonMinifier.String("application/json", jsonBody)
	assert.Nil(t, err)
	return minified
}

func generateHandler(t *testing.T, expectedBody string, expectedMethod string, responseBody string, responseHeader int) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		t.Helper()

		actualBody, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, minifyJSON(t, expectedBody), string(actualBody), "expected and actual bodies should be equal")
		assert.True(t, req.Method == expectedMethod)

		res.WriteHeader(responseHeader)
		fmt.Fprintf(res, responseBody)
	}
}

func generateHeadHandler(t *testing.T, responseHeader int) http.HandlerFunc {
	return generateHandler(t, "", http.MethodHead, "", responseHeader)
}

func generateGetHandler(t *testing.T, responseBody string, responseHeader int) http.HandlerFunc {
	return generateHandler(t, "", http.MethodGet, responseBody, responseHeader)
}

func generatePostHandler(t *testing.T, expectedBody string, responseBody string, responseHeader int) http.HandlerFunc {
	return generateHandler(t, expectedBody, http.MethodPost, responseBody, responseHeader)
}

func generatePatchHandler(t *testing.T, expectedBody string, responseBody string, responseHeader int) http.HandlerFunc {
	return generateHandler(t, expectedBody, http.MethodPatch, responseBody, responseHeader)
}

func generateDeleteHandler(t *testing.T, responseBody string, responseHeader int) http.HandlerFunc {
	return generateHandler(t, "", http.MethodDelete, responseBody, responseHeader)
}
