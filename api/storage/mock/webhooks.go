package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"
)

func (m *MockDB) WebhookExists(db storage.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetWebhook(db storage.Querier, id uint64) (*models.Webhook, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.Webhook), args.Error(1)
}

func (m *MockDB) GetWebhookList(db storage.Querier, qf *models.QueryFilter) ([]models.Webhook, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.Webhook), args.Error(1)
}

func (m *MockDB) GetWebhookCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateWebhook(db storage.Querier, nu *models.Webhook) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateWebhook(db storage.Querier, updated *models.Webhook) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteWebhook(db storage.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) GetWebhooksByEventType(db storage.Querier, eventType string) ([]models.Webhook, error) {
	args := m.Called(db, eventType)
	return args.Get(0).([]models.Webhook), args.Error(1)
}
