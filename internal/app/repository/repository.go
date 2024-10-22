package repository

type URL struct {
	UUID      string `json:"uuid"`
	LongURL   string `json:"long_url"`
	ShortCode string `json:"short_code"`
}

type URLRepository interface {
	Set(url URL) error
	Get(shortCode string) (*URL, bool)
}
