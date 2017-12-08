package main

import (
	"testing"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/mock"
)

type mockWebhookExecutor struct {
	mock.Mock
}

var _ WebhookExecutor = (*mockWebhookExecutor)(nil)

func (m *mockWebhookExecutor) CallWebhook(wh models.Webhook, object interface{}, db storage.Querier, client storage.Storer) {
	m.Called(wh, object, db, client)
}

func TestCallWebhook(t *testing.T) {

}
