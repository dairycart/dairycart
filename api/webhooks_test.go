package main

import (
	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/mock"
)

type mockWebhookExecutor struct {
	mock.Mock
}

var _ WebhookExecutor = (*mockWebhookExecutor)(nil)

func (m *mockWebhookExecutor) CallWebhook(wh models.Webhook, object interface{}, client storage.Storer) error {
	args := m.Called(wh, object, client)
	return args.Error(0)
}
