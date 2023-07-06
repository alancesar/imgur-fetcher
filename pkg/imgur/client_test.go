package imgur

import (
	"github.com/alancesar/imgur-fetcher/pkg/imgur/testdata"
	"net/http"
	"reflect"
	"testing"
)

func TestRestClient_GetMedia(t *testing.T) {
	type fields struct {
		httpClient *http.Client
	}
	type args struct {
		imageID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Media
		wantErr bool
	}{
		{
			name: "Should return image properly",
			fields: fields{
				httpClient: testdata.NewHTTPClient([]byte(testdata.ImgurImageResponse), http.StatusOK, nil),
			},
			args: args{
				imageID: "some-image-id",
			},
			want: Media{
				ID:          "some-image-id",
				Title:       "Some image title",
				Description: "Some image description",
				Link:        "https://i.imgur.com/some-image.jpg",
				Type:        "image/jpeg",
			},
			wantErr: false,
		},
		{
			name: "Should return video properly",
			fields: fields{
				httpClient: testdata.NewHTTPClient([]byte(testdata.ImgurVideoResponse), http.StatusOK, nil),
			},
			args: args{
				imageID: "some-video-id",
			},
			want: Media{
				ID:          "some-video-id",
				Title:       "Some video title",
				Description: "Some video description",
				Link:        "https://i.imgur.com/some-video.mp4",
				Type:        "video/mp4",
				MP4:         "https://i.imgur.com/some-video.mp4",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.httpClient)
			got, err := c.GetMedia(tt.args.imageID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMedia() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMedia() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImgurClient_GetAlbum(t *testing.T) {
	type fields struct {
		httpClient *http.Client
	}
	type args struct {
		albumID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    Album
		wantErr bool
	}{
		{
			name: "Should get images from album properly",
			fields: fields{
				httpClient: testdata.NewHTTPClient([]byte(testdata.ImgurAlbumResponse), http.StatusOK, nil),
			},
			args: args{
				albumID: "some-album-id",
			},
			want: Album{
				ID:          "some-album",
				Title:       "Some album title",
				Description: "Some album description",
				Link:        "https://imgur.com/a/some-album",
				Images: []Media{
					{
						ID:          "some-image-id-1",
						Title:       "Some image #1 title",
						Description: "Some image #1 description",
						Link:        "https://i.imgur.com/some-image-1.jpg",
						Type:        "image/jpeg",
					},
					{
						ID:          "some-image-id-2",
						Title:       "Some image #2 title",
						Description: "Some image #2 description",
						Link:        "https://i.imgur.com/some-image-2.jpg",
						Type:        "image/jpeg",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewClient(tt.fields.httpClient)
			got, err := c.GetAlbum(tt.args.albumID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAlbum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAlbum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImgurImage_GetHigherQualityImageURL(t *testing.T) {
	type fields struct {
		Type string
		Link string
		MP4  string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "Should return the value from link if it the type is image/jpeg",
			fields: fields{
				Type: "image/jpeg",
				Link: "https://link",
			},
			want: "https://link",
		},
		{
			name: "Should return the value from link if it the type is image/gif and mp4 value is empty",
			fields: fields{
				Type: "image/gif",
				Link: "https://link",
			},
			want: "https://link",
		},
		{
			name: "Should return the value from mp4 if it is present and the type is image/gif",
			fields: fields{
				Type: "image/gif",
				Link: "https://link",
				MP4:  "https://mp4",
			},
			want: "https://mp4",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Media{
				Link: tt.fields.Link,
				MP4:  tt.fields.MP4,
				Type: tt.fields.Type,
			}
			if got := m.HigherQualityURL(); got != tt.want {
				t.Errorf("HigherQualityURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetMediaByURL(t *testing.T) {
	type fields struct {
		httpClient *http.Client
	}
	type args struct {
		rawURL string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Media
		wantErr bool
	}{
		{
			name: "Should return media from URL",
			fields: fields{
				httpClient: testdata.NewHTTPClient([]byte(testdata.ImgurAlbumResponse), http.StatusOK, nil),
			},
			args: args{
				rawURL: "https://imgur.com/a/some-album",
			},
			want: []Media{
				{
					ID:          "some-image-id-1",
					Title:       "Some image #1 title",
					Description: "Some image #1 description",
					Link:        "https://i.imgur.com/some-image-1.jpg",
					Type:        "image/jpeg",
				},
				{
					ID:          "some-image-id-2",
					Title:       "Some image #2 title",
					Description: "Some image #2 description",
					Link:        "https://i.imgur.com/some-image-2.jpg",
					Type:        "image/jpeg",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Client{
				httpClient: tt.fields.httpClient,
			}
			got, err := c.GetMediaByURL(tt.args.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetMediaByURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetMediaByURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
