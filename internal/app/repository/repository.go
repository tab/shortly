package repository

type URL struct {
	LongURL   string
	ShortCode string
}

type URLRepository interface {
	Set(url URL)
	Get(shortCode string) (*URL, bool)
}
