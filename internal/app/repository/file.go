package repository

import (
	"encoding/json"
	"io"
	"os"

	"shortly/internal/app/errors"
)

type File interface {
	Load() (*Memento, error)
	Save(m *Memento) error
}

type fileRepo struct {
	filePath string
}

func NewFileRepository(filePath string) File {
	return &fileRepo{filePath: filePath}
}

func (f *fileRepo) Load() (*Memento, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, errors.ErrFailedToOpenFile
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	memento := &Memento{State: []URL{}}

	for {
		var url URL
		if err = decoder.Decode(&url); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, errors.ErrorFailedToReadFromFile
		}
		memento.State = append(memento.State, url)
	}

	return memento, nil
}

func (f *fileRepo) Save(memento *Memento) error {
	file, err := os.OpenFile(f.filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return errors.ErrFailedToOpenFile
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, row := range memento.State {
		if err = encoder.Encode(row); err != nil {
			return errors.ErrFailedToWriteToFile
		}
	}

	return nil
}
