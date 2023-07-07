package controller

import (
	"context"
	"encoding/json"
	"github.com/alancesar/imgur-fetcher/pkg/imgur"
	"github.com/alancesar/imgur-fetcher/pkg/media"
	"net/http"
)

type (
	Client interface {
		GetMediaByURL(url string) ([]imgur.Media, error)
	}

	Publisher interface {
		Publish(ctx context.Context, req media.Media) error
	}

	Controller struct {
		httpClient *http.Client
		client     Client
		publisher  Publisher
	}

	Response struct {
		URLs []string `json:"urls"`
	}
)

func New(httpClient *http.Client, client Client, publisher Publisher) *Controller {
	return &Controller{
		httpClient: httpClient,
		client:     client,
		publisher:  publisher,
	}
}

func (c Controller) GetMediaByURL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m, err := c.client.GetMediaByURL(req.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response Response
	response.URLs = make([]string, len(m), len(m))
	for i, m := range m {
		response.URLs[i] = m.HigherQualityURL()
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c Controller) PublishMedia(w http.ResponseWriter, r *http.Request) {
	var m media.Media
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	res, err := c.httpClient.Head(m.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if res.StatusCode >= http.StatusBadRequest {
		w.WriteHeader(res.StatusCode)
		return
	}

	m.URL = res.Request.URL.String()
	if err := c.publisher.Publish(r.Context(), m); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
