package transport

import (
	"context"
	"github.com/alancesar/imgur-fetcher/pkg/transport/testdata"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

var (
	parsedURL, _ = url.Parse("https://localhost")
)

func TestUserAgentRoundTripper_RoundTrip(t *testing.T) {
	type fields struct {
		userAgent string
		next      http.RoundTripper
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Response
		wantErr bool
	}{
		{
			name: "Should add User-Agent header properly",
			fields: fields{
				userAgent: "some-user-agent",
				next:      &testdata.FakedRoundTripper{},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					URL:    parsedURL,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body:          io.NopCloser(strings.NewReader(testdata.SampleRequestBody)),
					ContentLength: int64(len(testdata.SampleRequestBody)),
				},
			},
			want: &http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(strings.NewReader(testdata.SampleResponseBody)),
				ContentLength: int64(len(testdata.SampleResponseBody)),
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    parsedURL,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
						"User-Agent":   []string{"some-user-agent"},
					},
					Body:          io.NopCloser(strings.NewReader(testdata.SampleRequestBody)),
					ContentLength: int64(len(testdata.SampleRequestBody)),
				},
			},
			wantErr: false,
		},
		{
			name: "Should add User-Agent header properly even the request body is nil",
			fields: fields{
				userAgent: "some-user-agent",
				next:      &testdata.FakedRoundTripper{},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					URL:    parsedURL,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body:          nil,
					ContentLength: 0,
				},
			},
			want: &http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(strings.NewReader(testdata.SampleResponseBody)),
				ContentLength: int64(len(testdata.SampleResponseBody)),
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    parsedURL,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
						"User-Agent":   []string{"some-user-agent"},
					},
					Body:          nil,
					ContentLength: 0,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			roundTripper := NewUserAgentRoundTripper(tt.fields.userAgent, tt.fields.next)
			got, err := roundTripper.RoundTrip(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoundTrip() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthorizationRoundTripper_RoundTrip(t *testing.T) {
	type fields struct {
		provider      TokenProvider
		next          http.RoundTripper
		authorization string
	}
	type args struct {
		r *http.Request
	}
	tests := []struct {
		name            string
		fields          fields
		args            args
		want            *http.Response
		wantErr         bool
		wantStoredToken string
	}{
		{
			name: "Should bind with token from provider",
			fields: fields{
				provider: func(_ context.Context) (string, error) {
					return "some-token", nil
				},
				next: &testdata.FakedRoundTripper{},
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					URL:    parsedURL,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body:          io.NopCloser(strings.NewReader(testdata.SampleRequestBody)),
					ContentLength: int64(len(testdata.SampleRequestBody)),
				},
			},
			want: &http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(strings.NewReader(testdata.SampleResponseBody)),
				ContentLength: int64(len(testdata.SampleResponseBody)),
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    parsedURL,
					Header: http.Header{
						"Content-Type":  []string{"application/json"},
						"Authorization": []string{"some-token"},
					},
					Body:          io.NopCloser(strings.NewReader(testdata.SampleRequestBody)),
					ContentLength: int64(len(testdata.SampleRequestBody)),
				},
			},
			wantStoredToken: "some-token",
			wantErr:         false,
		},
		{
			name: "Should use the stored token",
			fields: fields{
				provider: func(_ context.Context) (string, error) {
					return "fresh-token", nil
				},
				next:          &testdata.FakedRoundTripper{},
				authorization: "stored-token",
			},
			args: args{
				r: &http.Request{
					Method: http.MethodPost,
					URL:    parsedURL,
					Header: http.Header{
						"Content-Type": []string{"application/json"},
					},
					Body:          io.NopCloser(strings.NewReader(testdata.SampleRequestBody)),
					ContentLength: int64(len(testdata.SampleRequestBody)),
				},
			},
			want: &http.Response{
				Status:        http.StatusText(http.StatusOK),
				StatusCode:    http.StatusOK,
				Body:          io.NopCloser(strings.NewReader(testdata.SampleResponseBody)),
				ContentLength: int64(len(testdata.SampleResponseBody)),
				Request: &http.Request{
					Method: http.MethodPost,
					URL:    parsedURL,
					Header: http.Header{
						"Content-Type":  []string{"application/json"},
						"Authorization": []string{"stored-token"},
					},
					Body:          io.NopCloser(strings.NewReader(testdata.SampleRequestBody)),
					ContentLength: int64(len(testdata.SampleRequestBody)),
				},
			},
			wantStoredToken: "some-token",
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &AuthorizationRoundTripper{
				provider:      tt.fields.provider,
				next:          tt.fields.next,
				authorization: tt.fields.authorization,
			}
			got, err := p.RoundTrip(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("RoundTrip() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RoundTrip() got = %v, want %v", got, tt.want)
			}
		})
	}
}
