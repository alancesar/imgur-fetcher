package testdata

import (
	"bytes"
	"io"
	"net/http"
)

type (
	Transport struct {
		response   []byte
		statusCode int
		err        error
	}
)

func (m Transport) RoundTrip(request *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:        http.StatusText(m.statusCode),
		StatusCode:    m.statusCode,
		Body:          io.NopCloser(bytes.NewReader(m.response)),
		ContentLength: int64(len(m.response)),
		Request:       request,
	}, m.err
}

func NewHTTPClient(response []byte, statusCode int, err error) *http.Client {
	transport := Transport{
		response:   response,
		statusCode: statusCode,
		err:        err,
	}

	return &http.Client{
		Transport: transport,
	}
}
