package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/mock"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockWebhookExecutor struct {
	mock.Mock
}

var _ WebhookExecutor = (*mockWebhookExecutor)(nil)

func (m *mockWebhookExecutor) CallWebhook(wh models.Webhook, object interface{}, db storage.Querier, client storage.Storer) {
	m.Called(wh, object, db, client)
}

type testBreakableStruct struct {
	Thing json.Number `json:"thing"`
}

func handlerGenerator(handlers map[string]http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for path, handlerFunc := range handlers {
			if r.URL.Path == path {
				handlerFunc(w, r)
				return
			}
		}
	})
}

func generateHandler(t *testing.T, expectedBody string, responseBody string, responseHeader int) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		actualBody, err := ioutil.ReadAll(req.Body)
		require.Nil(t, err)
		assert.Equal(t, expectedBody, string(actualBody), "expected and actual bodies should be equal")

		assert.True(t, req.Method == http.MethodPost)

		res.WriteHeader(responseHeader)
		fmt.Fprintf(res, responseBody)
	}
}

func TestCallWebhook(t *testing.T) {
	t.Run("normal operation", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", "{}", http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{URL: ts.URL}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), nil)

		whe.CallWebhook(exampleWebhook, &models.Product{}, db, client)
	})

	t.Run("with xml", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", "{}", http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{
			ContentType: "application/xml",
			URL:         ts.URL,
		}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), nil)

		whe.CallWebhook(exampleWebhook, &models.Product{}, db, client)
	})

	t.Run("with invalid object", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", "{}", http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{URL: ts.URL}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), nil)

		whe.CallWebhook(exampleWebhook, &testBreakableStruct{Thing: "broken"}, db, client)
	})

	t.Run("with invalid URL", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", "{}", http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{URL: ":"}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), nil)

		whe.CallWebhook(exampleWebhook, &models.Product{}, db, client)
	})

	t.Run("with error", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", "{}", http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{URL: ts.URL}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), generateExampleTimeForTests(), generateArbitraryError())

		whe.CallWebhook(exampleWebhook, &models.Product{}, db, client)
	})
}
