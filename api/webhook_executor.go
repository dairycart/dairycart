package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"
)

type WebhookExecutor interface {
	CallWebhook(models.Webhook, interface{}, storage.Querier, storage.Storer)
}

type webhookExecutor struct{}

var _ WebhookExecutor = (*webhookExecutor)(nil)

func (whe *webhookExecutor) CallWebhook(wh models.Webhook, object interface{}, db storage.Querier, client storage.Storer) {
	var body []byte
	var err error

	wel := &models.WebhookExecutionLog{WebhookID: wh.ID}

	switch strings.ToLower(wh.ContentType) {
	case "application/json":
		body, err = json.Marshal(object)
		if err != nil {
			log.Printf("error encountered executing webhook: %v", err)
			return
		}
	case "application/xml":
		body, err = xml.Marshal(object)
		if err != nil {
			log.Printf("error encountered executing webhook: %v", err)
			return
		}
	default:
		log.Printf("invalid content type: %s", wh.ContentType)
		return
	}

	res, err := http.Post(wh.URL, wh.ContentType, bytes.NewBuffer(body))
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
