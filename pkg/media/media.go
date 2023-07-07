package media

import "path"

type (
	Media struct {
		URL      string   `json:"url"`
		Filename string   `json:"filename"`
		Parent   []string `json:"parent"`
	}
)

func (m Media) Path() string {
	return path.Join(path.Join(m.Parent...), m.Filename)
}
