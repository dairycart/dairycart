package main

import (
	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/stretchr/testify/mock"
)

type mockWebhookExecutor struct {
	mock.Mock
}

var _ WebhookExecutor = (*mockWebhookExecutor)(nil)

func (m *mockWebhookExecutor) CallWebhook(wh models.Webhook, object interface{}) error {
	args := m.Called(wh, object)
	return args.Error(0)
}
