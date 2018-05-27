package api

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
	"strings"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

type WebhookExecutor interface {
	CallWebhook(models.Webhook, interface{}, database.Querier, database.Storer)
}

type webhookExecutor struct {
	*http.Client
}

var _ WebhookExecutor = (*webhookExecutor)(nil)

func (whe *webhookExecutor) CallWebhook(wh models.Webhook, object interface{}, db database.Querier, client database.Storer) {
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
