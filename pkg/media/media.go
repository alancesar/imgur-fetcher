package media

type (
	Media struct {
		URL    string   `json:"url"`
		Parent []string `json:"parent"`
	}
)
