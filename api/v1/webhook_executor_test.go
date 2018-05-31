package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
	"github.com/dairycart/dairycart/storage/v1/database/mock"

	"github.com/stretchr/testify/mock"
)

var emptyJSONObj = []byte("{}")

type mockWebhookExecutor struct{}

var _ WebhookExecutor = (*mockWebhookExecutor)(nil)

func (m *mockWebhookExecutor) CallWebhook(wh models.Webhook, object interface{}, db database.Querier, client database.Storer) {
}

type testBreakableStruct struct {
	Thing json.Number `json:"thing"`
}

func TestCallWebhook(t *testing.T) {
	t.Run("normal operation", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", emptyJSONObj, http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{URL: ts.URL}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), buildTestTime(), nil)

		whe.CallWebhook(exampleWebhook, &models.Product{}, db, client)
	})

	t.Run("with xml", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", emptyJSONObj, http.StatusOK),
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
			Return(uint64(0), buildTestTime(), nil)

		whe.CallWebhook(exampleWebhook, &models.Product{}, db, client)
	})

	t.Run("with invalid object", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", emptyJSONObj, http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{URL: ts.URL}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), buildTestTime(), nil)

		whe.CallWebhook(exampleWebhook, &testBreakableStruct{Thing: "broken"}, db, client)
	})

	t.Run("with invalid URL", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", emptyJSONObj, http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{URL: ":"}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), buildTestTime(), nil)

		whe.CallWebhook(exampleWebhook, &models.Product{}, db, client)
	})

	t.Run("with error", func(*testing.T) {
		db := setupTestVariablesWithMock(t).PlainDB

		handlers := map[string]http.HandlerFunc{
			"": generateHandler(t, "", emptyJSONObj, http.StatusOK),
		}
		ts := httptest.NewServer(handlerGenerator(handlers))
		defer ts.Close()

		exampleWebhook := models.Webhook{URL: ts.URL}
		whe := &webhookExecutor{Client: ts.Client()}

		client := &dairymock.MockDB{}
		client.On("CreateWebhookExecutionLog", mock.Anything, mock.Anything).
			Return(uint64(0), buildTestTime(), generateArbitraryError())

		whe.CallWebhook(exampleWebhook, &models.Product{}, db, client)
	})
}
