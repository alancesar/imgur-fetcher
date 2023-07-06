package testdata

import (
	"io"
	"net/http"
	"strings"
)

const (
	SampleRequestBody  = `{"request": "foo"}`
	SampleResponseBody = `{"response": "bar"}`
)

type (
	FakedRoundTripper struct{}
)

func (m FakedRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:        http.StatusText(http.StatusOK),
		StatusCode:    http.StatusOK,
		Body:          io.NopCloser(strings.NewReader(SampleResponseBody)),
		ContentLength: int64(len(SampleResponseBody)),
		Request:       r,
	}, nil
}
