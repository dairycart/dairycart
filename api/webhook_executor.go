package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"
)

type WebhookExecutor interface {
	CallWebhook(models.Webhook, interface{}, storage.Querier, storage.Storer)
}

type webhookExecutor struct {
	*http.Client
}

var _ WebhookExecutor = (*webhookExecutor)(nil)

func (whe *webhookExecutor) CallWebhook(wh models.Webhook, object interface{}, db storage.Querier, client storage.Storer) {
	var (
		body        []byte
		err         error
		marshalFunc func(v interface{}) ([]byte, error)
	)

	wel := &models.WebhookExecutionLog{WebhookID: wh.ID}

	switch strings.ToLower(wh.ContentType) {
	case "application/xml":
		marshalFunc = xml.Marshal
	default:
		marshalFunc = json.Marshal
	}

	body, err = marshalFunc(object)
	if err != nil {
		log.Printf("error encountered executing webhook: %v", err)
		return
	}

	res, err := whe.Client.Post(wh.URL, wh.ContentType, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error encountered executing webhook: %v", err)
		return
	}

	if res != nil {
		wel.StatusCode = res.StatusCode
	}
	if wel.StatusCode >= http.StatusOK && wel.StatusCode <= http.StatusMultipleChoices {
		wel.Succeeded = true
	}

	_, _, err = client.CreateWebhookExecutionLog(db, wel)
	if err != nil {
		log.Printf("error encountered logging webhook execution: %v", err)
	}
}
