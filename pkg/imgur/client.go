package imgur

import (
	"encoding/json"
	"fmt"
	"github.com/alancesar/imgur-fetcher/pkg/status"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
)

const (
	apiPath      = "https://api.imgur.com/3"
	gifImageType = "image/gif"
)

type (
	Client struct {
		httpClient *http.Client
	}

	Request struct {
		ID      string
		IsAlbum bool
	}

	Response[T any] struct {
		Data T `json:"data"`
	}

	Media struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Link        string `json:"link"`
		Type        string `json:"type"`
		MP4         string `json:"mp4"`
	}

	Album struct {
		ID          string  `json:"id"`
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Link        string  `json:"link"`
		Images      []Media `json:"images"`
	}
)

func (m Media) HigherQualityURL() string {
	if m.Type == gifImageType && m.MP4 != "" {
		return m.MP4
	}

	return m.Link
}

func NewClient(httpClient *http.Client) *Client {
	return &Client{
		httpClient: httpClient,
	}
}

func (c Client) GetMediaByURL(rawURL string) ([]Media, error) {
	request, err := compileURL(rawURL)
	if err != nil {
		return nil, err
	}

	if request.IsAlbum {
		album, err := c.GetAlbum(request.ID)
		if err != nil {
			return nil, err
		}

		return album.Images, nil
	}

	media, err := c.GetMedia(request.ID)
	return []Media{media}, nil
}

func (c Client) GetMedia(imageID string) (Media, error) {
	var output Response[Media]
	formattedURL := fmt.Sprintf("%s/image/%s", apiPath, imageID)
	err := c.doGet(formattedURL, &output)
	return output.Data, err
}

func (c Client) GetAlbum(albumID string) (Album, error) {
	var output Response[Album]
	formattedURL := fmt.Sprintf("%s/album/%s", apiPath, albumID)
	err := c.doGet(formattedURL, &output)
	return output.Data, err
}

func (c Client) doGet(url string, output any) error {
	res, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("%w: %s", status.ErrNotFound, url)
	} else if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("%w: %d (%s): %s", status.ErrBadStatus, res.StatusCode, res.Status, url)
	}

	return json.NewDecoder(res.Body).Decode(&output)
}

func compileURL(rawURL string) (Request, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return Request{}, err
	}

	isAlbum := isAlbumURL(parsedURL.Path)
	id := path.Base(parsedURL.Path)

	if ext := filepath.Ext(id); ext != "" {
		id = strings.ReplaceAll(id, ext, "")
		return Request{
			ID:      id,
			IsAlbum: isAlbum,
		}, nil
	}

	return Request{
		ID:      id,
		IsAlbum: isAlbum,
	}, nil
}

func isAlbumURL(path string) bool {
	return strings.Contains(path, "/a/")
}
