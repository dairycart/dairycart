package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

////////////////////////////////////////////////////////
//                                                    //
//                 HTTP Handler Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestWebhookRetrievalHandler(t *testing.T) {
	exampleWebhook := &models.Webhook{}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/webhook/1", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent webhook", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/webhook/1", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving webhook", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/webhook/1", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestWebhookListHandler(t *testing.T) {
	exampleWebhook := models.Webhook{}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhookCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetWebhookList", mock.Anything, mock.Anything).
			Return([]models.Webhook{exampleWebhook}, nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/webhooks", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with error retrieving webhook count", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhookCount", mock.Anything, mock.Anything).
			Return(uint64(3), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/webhooks", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error retrieving webhook list", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhookCount", mock.Anything, mock.Anything).
			Return(uint64(3), nil)
		testUtil.MockDB.On("GetWebhookList", mock.Anything, mock.Anything).
			Return([]models.Webhook{}, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodGet, "/v1/webhooks", nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestWebhookCreationHandler(t *testing.T) {
	exampleWebhookCreationInput := `
		{

		}
	`

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("CreateWebhook", mock.Anything, mock.Anything).
			Return(uint64(1), buildTestTime(), nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/webhook", strings.NewReader(exampleWebhookCreationInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusCreated)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/webhook", strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with error creating webhook", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("CreateWebhook", mock.Anything, mock.Anything).
			Return(uint64(1), buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodPost, "/v1/webhook", strings.NewReader(exampleWebhookCreationInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestWebhookDeletionHandler(t *testing.T) {
	exampleWebhook := &models.Webhook{}

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, nil)
		testUtil.MockDB.On("DeleteWebhook", mock.Anything, exampleWebhook.ID).
			Return(buildTestTime(), nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/webhook/%d", exampleWebhook.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with nonexistent webhook", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/webhook/%d", exampleWebhook.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving webhook", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/webhook/%d", exampleWebhook.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error deleting webhook", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, nil)
		testUtil.MockDB.On("DeleteWebhook", mock.Anything, exampleWebhook.ID).
			Return(buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v1/webhook/%d", exampleWebhook.ID), nil)
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}

func TestWebhookUpdateHandler(t *testing.T) {
	exampleWebhook := &models.Webhook{}

	exampleWebhookUpdateInput := `
		{

		}
	`

	t.Run("optimal conditions", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, nil)
		testUtil.MockDB.On("UpdateWebhook", mock.Anything, mock.Anything).
			Return(buildTestTime(), nil)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/webhook/1", strings.NewReader(exampleWebhookUpdateInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusOK)
	})

	t.Run("with invalid input", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/webhook/1", strings.NewReader(exampleGarbageInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusBadRequest)
	})

	t.Run("with nonexistent error", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, sql.ErrNoRows)
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/webhook/1", strings.NewReader(exampleWebhookUpdateInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusNotFound)
	})

	t.Run("with error retrieving webhook", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/webhook/1", strings.NewReader(exampleWebhookUpdateInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})

	t.Run("with error updating webhook", func(*testing.T) {
		testUtil := setupTestVariablesWithMock(t)
		testUtil.MockDB.On("GetWebhook", mock.Anything, exampleWebhook.ID).
			Return(exampleWebhook, nil)
		testUtil.MockDB.On("UpdateWebhook", mock.Anything, mock.Anything).
			Return(buildTestTime(), generateArbitraryError())
		config := buildServerConfigFromTestUtil(testUtil)
		SetupAPIRouter(config)

		req, err := http.NewRequest(http.MethodPatch, "/v1/webhook/1", strings.NewReader(exampleWebhookUpdateInput))
		assert.NoError(t, err)

		testUtil.Router.ServeHTTP(testUtil.Response, req)
		assertStatusCode(t, testUtil, http.StatusInternalServerError)
	})
}
