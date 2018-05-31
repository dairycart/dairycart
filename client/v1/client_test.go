package dairyclient_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dairycart/dairycart/client/v1"

	"github.com/stretchr/testify/assert"
)

////////////////////////////////////////////////////////
//                                                    //
//                 Constructor Tests                  //
//                                                    //
////////////////////////////////////////////////////////

func TestNewV1Client(t *testing.T) {

	t.Run("normal usage", func(t *testing.T) {
		ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
		defer ts.Close()
		c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

		assert.NotNil(t, c.AuthCookie)
		assert.Nil(t, err)
	})

	t.Run("invalid URL", func(t *testing.T) {
		c, err := dairyclient.NewV1Client(":", exampleUsername, examplePassword, nil)
		assert.Nil(t, c)
		assert.NotNil(t, err)
	})

	t.Run("login failure", func(t *testing.T) {
		ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
		ts.Close()
		c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

		assert.Nil(t, c)
		assert.NotNil(t, err)
	})

	t.Run("without cookie", func(t *testing.T) {
		ts := httptest.NewTLSServer(obligatoryLoginHandler(false))
		defer ts.Close()
		c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

		assert.Nil(t, c)
		assert.NotNil(t, err)
	})
}

func TestNewV1ClientFromCookie(t *testing.T) {

	t.Run("normal use", func(*testing.T) {
		ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
		defer ts.Close()

		_, err := dairyclient.NewV1ClientFromCookie(ts.URL, &http.Cookie{}, ts.Client())
		assert.Nil(t, err)
	})

	t.Run("with invalid URL", func(*testing.T) {
		_, err := dairyclient.NewV1ClientFromCookie(":", &http.Cookie{}, http.DefaultClient)
		assert.NotNil(t, err)
	})
}

func TestBuildURL(t *testing.T) {

	ts := httptest.NewTLSServer(http.NotFoundHandler())
	defer ts.Close()
	c := buildTestClient(t, ts)

	t.Run("normal usage", func(t *testing.T) {
		expected := fmt.Sprintf("%s/v1/things/stuff?query=params", ts.URL)
		exampleParams := map[string]string{
			"query": "params",
		}
		actual, err := c.BuildURL(exampleParams, "things", "stuff")

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "BuildURL doesn't return the correct result. Expected `%s`, got `%s`", expected, actual)
	})

	t.Run("invalid URL", func(t *testing.T) {
		actual, err := c.BuildURL(nil, `%gh&%ij`)

		assert.NotNil(t, err)
		assert.Empty(t, actual)
	})
}
