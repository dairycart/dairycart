package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) WebhookExecutionLogExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetWebhookExecutionLog(db database.Querier, id uint64) (*models.WebhookExecutionLog, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.WebhookExecutionLog), args.Error(1)
}

func (m *MockDB) GetWebhookExecutionLogList(db database.Querier, qf *models.QueryFilter) ([]models.WebhookExecutionLog, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.WebhookExecutionLog), args.Error(1)
}

func (m *MockDB) GetWebhookExecutionLogCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateWebhookExecutionLog(db database.Querier, nu *models.WebhookExecutionLog) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateWebhookExecutionLog(db database.Querier, updated *models.WebhookExecutionLog) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteWebhookExecutionLog(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
