package controller

import (
	"encoding/json"
	"github.com/alancesar/imgur-fetcher/pkg/imgur"
	"net/http"
)

type (
	Client interface {
		GetMediaByURL(url string) ([]imgur.Media, error)
	}

	Controller struct {
		client Client
	}

	Response struct {
		URLs []string `json:"urls"`
	}
)

func New(client Client) *Controller {
	return &Controller{
		client: client,
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
