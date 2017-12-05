package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/dairycart/dairycart/api/storage/models"

	"github.com/pkg/errors"
)

type WebhookExecutor interface {
	CallWebhook(models.Webhook, interface{}) error
}

type webhookExecutor struct{}

var _ WebhookExecutor = (*webhookExecutor)(nil)

func (whe *webhookExecutor) CallWebhook(wh models.Webhook, object interface{}) error {
	var body []byte
	var err error

	switch strings.ToLower(wh.ContentType) {
	case "application/json":
		body, err = json.Marshal(object)
		if err != nil {
			return err
		}
	case "application/xml":
		body, err = xml.Marshal(object)
		if err != nil {
			return err
		}
	default:
		return errors.New(fmt.Sprintf("invalid content type: %s", wh.ContentType))
	}

	http.Post(wh.URL, wh.ContentType, bytes.NewBuffer(body))
	return nil
}
