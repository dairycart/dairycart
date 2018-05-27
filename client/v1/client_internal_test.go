// +build !exported

package dairyclient

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	exampleID       = 666
	exampleURL      = `http://www.dairycart.com`
	exampleUsername = `username`
	examplePassword = `password` // lol not really
	exampleSKU      = `sku`
	exampleBadJSON  = `{"invalid lol}`
)

////////////////////////////////////////////////////////
//                                                    //
//                 Helper Functions                   //
//                                                    //
////////////////////////////////////////////////////////

func handlerGenerator(handlers map[string]func(res http.ResponseWriter, req *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for path, handlerFunc := range handlers {
			if r.URL.Path == path {
				handlerFunc(w, r)
				return
			}
		}
	})
}

func createInternalClient(t *testing.T, ts *httptest.Server) *V1Client {
	u, err := url.Parse(ts.URL)
	assert.Nil(t, err, "no error should be returned when parsing a test server's URL")

	c := &V1Client{
		Client: ts.Client(),
		AuthCookie: &http.Cookie{
			Name: "dairycart",
		},
		URL: u,
	}
	return c
}

////////////////////////////////////////////////////////
//                                                    //
//                   Actual Tests                     //
//                                                    //
////////////////////////////////////////////////////////

func TestExecuteRequestAddsCookieToRequests(t *testing.T) {

	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			cookies := req.Cookies()
			if len(cookies) == 0 {
				assert.FailNow(t, "no cookies attached to the request")
			}

			cookieFound := false
			for _, c := range cookies {
				if c.Name == "dairycart" {
					cookieFound = true
				}
			}
			assert.True(t, cookieFound)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", ts.URL, exampleEndpoint), nil)
	assert.Nil(t, err, "no error should be returned when creating a new request")

	c.executeRequest(req)
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestUnexportedBuildURL(t *testing.T) {
	ts := httptest.NewTLSServer(http.NotFoundHandler())
	defer ts.Close()
	c := createInternalClient(t, ts)

	testCases := []struct {
		query    map[string]string
		parts    []string
		expected string
	}{
		{
			query:    nil,
			parts:    []string{""},
			expected: fmt.Sprintf("%s/v1/", ts.URL),
		},
		{
			query:    nil,
			parts:    []string{"things", "and", "stuff"},
			expected: fmt.Sprintf("%s/v1/things/and/stuff", ts.URL),
		},
		{
			query:    map[string]string{"param": "value"},
			parts:    []string{"example"},
			expected: fmt.Sprintf("%s/v1/example?param=value", ts.URL),
		},
	}

	for _, tc := range testCases {
		actual := c.buildURL(tc.query, tc.parts...)
		assert.Equal(t, tc.expected, actual, "expected and actual built URLs don't match")
	}
}

func TestExists(t *testing.T) {

	var normalEndpointCalled bool
	var fourOhFourEndpointCalled bool

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		"/v1/normal": func(res http.ResponseWriter, req *http.Request) {
			normalEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodHead, "exists should be making HEAD requests")
			res.WriteHeader(http.StatusOK)
		},
		"/v1/four_oh_four": func(res http.ResponseWriter, req *http.Request) {
			fourOhFourEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodHead, "exists should be making HEAD requests")
			res.WriteHeader(http.StatusNotFound)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	c := createInternalClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		actual, err := c.exists(c.buildURL(nil, "normal"))
		assert.Nil(t, err)
		assert.True(t, actual, "exists should return false when the status code is %d", http.StatusOK)
		assert.True(t, normalEndpointCalled, "endpoint should have been called")
	})

	t.Run("not found", func(t *testing.T) {
		actual, err := c.exists(c.buildURL(nil, "four_oh_four"))
		assert.Nil(t, err)
		assert.False(t, actual, "exists should return false when the status code is %d", http.StatusNotFound)
		assert.True(t, fourOhFourEndpointCalled, "endpoint should have been called")
	})

	t.Run("failure executing request", func(t *testing.T) {
		ts.Close()
		actual, err := c.exists(c.buildURL(nil, "whatever"))
		assert.NotNil(t, err)
		assert.False(t, actual, "exists should return false when the status code is %d", http.StatusOK)
	})
}

func TestGet(t *testing.T) {

	var normalEndpointCalled bool
	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		"/v1/normal": func(res http.ResponseWriter, req *http.Request) {
			normalEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodGet, "get should be making GET requests")
			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		expected := struct {
			Things string `json:"things"`
		}{
			Things: "stuff",
		}

		actual := struct {
			Things string `json:"things"`
		}{}

		err := c.get(c.buildURL(nil, "normal"), &actual)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "actual struct should equal expected struct")
		assert.True(t, normalEndpointCalled, "endpoint should have been called")
	})

	t.Run("nil input", func(t *testing.T) {
		nilErr := c.get(c.buildURL(nil, "whatever"), nil)
		assert.NotNil(t, nilErr)
	})

	t.Run("non pointer input", func(t *testing.T) {
		actual := struct {
			Things string `json:"things"`
		}{}

		ptrErr := c.get(c.buildURL(nil, "whatever"), actual)
		assert.NotNil(t, ptrErr)
	})
}

func TestDelete(t *testing.T) {
	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		"/v1/normal": func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, http.MethodDelete, "delete should be making DELETE requests")
			fmt.Fprintf(res, "{}")
		},
		"/v1/five_hundred": func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, http.MethodDelete, "delete should be making DELETE requests")
			res.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(res, `
				{
					"status": 500,
					"message": "obligatory error"
				}
			`)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	c := createInternalClient(t, ts)

	t.Run("normal usage", func(t *testing.T) {
		err := c.delete(c.buildURL(nil, "normal"))
		assert.Nil(t, err)
	})

	t.Run("bad status code", func(t *testing.T) {
		u := c.buildURL(nil, "five_hundred")
		err := c.delete(u)
		assert.NotNil(t, err)
	})

	t.Run("failed request", func(t *testing.T) {
		ts.Close()
		err := c.delete(c.buildURL(nil, "whatever"))
		assert.NotNil(t, err)
	})
}

func TestMakeDataRequest(t *testing.T) {
	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		"/v1/whatever": func(res http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, http.MethodPost, "makeDataRequest should only be making PUT or POST requests")
			exampleResponse := `{"things":"stuff"}`

			bodyBytes, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			requestBody := string(bodyBytes)
			assert.Equal(t, requestBody, exampleResponse, "makeDataRequest should attach the correct JSON to the request body")

			fmt.Fprintf(res, exampleResponse)
		},
		"/v1/bad_json": func(res http.ResponseWriter, req *http.Request) {
			fmt.Fprintf(res, exampleBadJSON)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	c := createInternalClient(t, ts)

	t.Run("normal usage", func(*testing.T) {
		expected := struct {
			Things string `json:"things"`
		}{
			Things: "stuff",
		}

		actual := struct {
			Things string `json:"things"`
		}{}

		err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), expected, &actual)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "actual struct should equal expected struct")
	})

	t.Run("nil argument", func(*testing.T) {
		ptrErr := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), struct{}{}, struct{}{})
		assert.NotNil(t, ptrErr, "makeDataRequest should return an error when passed a non-pointer output param")
	})

	t.Run("non-pointer argument", func(*testing.T) {
		nilErr := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), struct{}{}, nil)
		assert.NotNil(t, nilErr, "makeDataRequest should return an error when passed a nil output param")
	})

	t.Run("invalid struct argument", func(*testing.T) {
		f := &testBreakableStruct{Thing: "dongs"}
		err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), f, &struct{}{})
		assert.NotNil(t, err, "makeDataRequest should return an error when passed an invalid input struct")
	})

	t.Run("unmarshal failure", func(*testing.T) {
		expected := struct {
			Things string `json:"things"`
		}{
			Things: "stuff",
		}

		actual := struct {
			Things string `json:"things"`
		}{}

		err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "bad_json"), expected, &actual)
		assert.NotNil(t, err)
	})

	t.Run("failed request", func(*testing.T) {
		expected := struct {
			Things string `json:"things"`
		}{
			Things: "stuff",
		}

		actual := struct {
			Things string `json:"things"`
		}{}

		ts.Close()
		err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), expected, &actual)
		assert.NotNil(t, err, "makeDataRequest should return an error when failing to execute request")
	})
}

func TestPost(t *testing.T) {

	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodPost, "post should only be making POST requests")
			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	expected := struct {
		Things string `json:"things"`
	}{
		Things: "stuff",
	}

	actual := struct {
		Things string `json:"things"`
	}{}

	exampleURI := c.buildURL(nil, "whatever")
	err := c.post(exampleURI, expected, &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "actual struct should equal expected struct")
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestPatch(t *testing.T) {

	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodPatch, "patch should only be making PATCH requests")
			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	expected := struct {
		Things string `json:"things"`
	}{
		Things: "stuff",
	}

	actual := struct {
		Things string `json:"things"`
	}{}

	err := c.patch(c.buildURL(nil, "whatever"), expected, &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "actual struct should equal expected struct")
	assert.True(t, endpointCalled, "endpoint should have been called")
}
