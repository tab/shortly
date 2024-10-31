package repository

type URL struct {
	UUID      string `json:"uuid"`
	LongURL   string `json:"long_url"`
	ShortCode string `json:"short_code"`
}

type Memento struct {
	State []URL `json:"state"`
}

type Repository interface {
	Set(url URL) error
	Get(shortCode string) (*URL, bool)
	CreateMemento() *Memento
	Restore(m *Memento)
}

func NewRepository() Repository {
	return NewInMemoryRepository()
}
