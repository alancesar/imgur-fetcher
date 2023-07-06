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

	Request struct {
		URL string `json:"url"`
	}

	Response struct {
		Urls []string `json:"urls"`
	}
)

func New(client Client) *Controller {
	return &Controller{
		client: client,
	}
}

func (c Controller) GetMediaByURL(w http.ResponseWriter, r *http.Request) {
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	media, err := c.client.GetMediaByURL(request.URL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var response Response
	response.Urls = make([]string, len(media), len(media))
	for i, m := range media {
		response.Urls[i] = m.HigherQualityURL()
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
