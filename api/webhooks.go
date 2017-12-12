package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"

	"github.com/go-chi/chi"
	"github.com/imdario/mergo"
)

func buildWebhookListRetrievalHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// WebhookListRetrievalHandler is a request handler that returns a list of Webhooks
	return func(res http.ResponseWriter, req *http.Request) {
		rawFilterParams := req.URL.Query()
		queryFilter := parseRawFilterParams(rawFilterParams)

		count, err := client.GetWebhookCount(db, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve count of webhooks from the database")
			return
		}

		webhooks, err := client.GetWebhookList(db, queryFilter)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve webhooks from the database")
			return
		}

		webhooksResponse := &ListResponse{
			Page:  queryFilter.Page,
			Limit: queryFilter.Limit,
			Count: count,
			Data:  webhooks,
		}
		json.NewEncoder(res).Encode(webhooksResponse)
	}
}

func buildWebhookListRetrievalByEventTypeHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// WebhookListRetrievalHandler is a request handler that returns a list of Webhooks
	return func(res http.ResponseWriter, req *http.Request) {
		eventType := chi.URLParam(req, "event_type")

		webhooks, err := client.GetWebhooksByEventType(db, eventType)
		if err != nil {
			notifyOfInternalIssue(res, err, "retrieve webhooks from the database")
			return
		}

		json.NewEncoder(res).Encode(webhooks)
	}
}

func buildWebhookCreationHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// WebhookCreationHandler is a request handler that creates a Webhook from user input
	return func(res http.ResponseWriter, req *http.Request) {
		newWebhook := &models.Webhook{}
		err := validateRequestInput(req, newWebhook)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		newID, createdOn, err := client.CreateWebhook(db, newWebhook)
		newWebhook.ID = newID
		newWebhook.CreatedOn = &models.Dairytime{Time: createdOn}
		if err != nil {
			notifyOfInternalIssue(res, err, "insert webhook into database")
			return
		}

		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(newWebhook)
	}
}

func buildWebhookDeletionHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// WebhookDeletionHandler is a request handler that deletes a single webhook
	return func(res http.ResponseWriter, req *http.Request) {
		webhookIDStr := chi.URLParam(req, "webhook_id")
		// eating this error because the router should have ensured this is an integer
		webhookID, _ := strconv.ParseUint(webhookIDStr, 10, 64)

		webhook, err := client.GetWebhook(db, webhookID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "webhook", webhookIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieving webhook from database")
			return
		}

		archivedOn, err := client.DeleteWebhook(db, webhookID)
		if err != nil {
			notifyOfInternalIssue(res, err, "archive webhook in database")
			return
		}
		webhook.ArchivedOn = &models.Dairytime{Time: archivedOn}

		json.NewEncoder(res).Encode(webhook)
	}
}

func buildWebhookUpdateHandler(db *sql.DB, client storage.Storer) http.HandlerFunc {
	// WebhookUpdateHandler is a request handler that can update webhooks
	return func(res http.ResponseWriter, req *http.Request) {
		webhookIDStr := chi.URLParam(req, "webhook_id")
		// eating this error because the router should have ensured this is an integer
		webhookID, _ := strconv.ParseUint(webhookIDStr, 10, 64)

		updatedWebhook := &models.Webhook{}
		err := validateRequestInput(req, updatedWebhook)
		if err != nil {
			notifyOfInvalidRequestBody(res, err)
			return
		}

		existingWebhook, err := client.GetWebhook(db, webhookID)
		if err == sql.ErrNoRows {
			respondThatRowDoesNotExist(req, res, "webhook", webhookIDStr)
			return
		} else if err != nil {
			notifyOfInternalIssue(res, err, "retrieve webhook from database")
			return
		}

		mergo.Merge(updatedWebhook, existingWebhook)

		updatedOn, err := client.UpdateWebhook(db, updatedWebhook)
		if err != nil {
			notifyOfInternalIssue(res, err, "update webhook in database")
			return
		}
		updatedWebhook.UpdatedOn = &models.Dairytime{Time: updatedOn}

		json.NewEncoder(res).Encode(updatedWebhook)
	}
}
