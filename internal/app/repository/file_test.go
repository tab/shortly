package repository

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"shortly/internal/app/errors"
)

func Test_FileStorageRepository_Load(t *testing.T) {
	type result struct {
		memento *Memento
		err     error
	}

	tests := []struct {
		name     string
		before   func(filePath string)
		expected result
	}{
		{
			name: "Success",
			before: func(filePath string) {
				file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				require.NoError(t, err)
				defer file.Close()

				_, err = file.WriteString(`{"uuid":"6455bd07-e431-4851-af3c-4f703f726639","long_url":"http://example.com","short_code":"abcd1234"}`)
				require.NoError(t, err)
			},
			expected: result{
				memento: &Memento{
					State: []URL{
						{
							UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
							LongURL:   "http://example.com",
							ShortCode: "abcd1234",
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "File not exists",
			before: func(filePath string) {
				os.Remove(filePath)
			},
			expected: result{
				memento: &Memento{
					State: []URL{},
				},
				err: nil,
			},
		},
		{
			name: "Multiple JSON objects",
			before: func(filePath string) {
				file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				require.NoError(t, err)
				defer file.Close()

				_, err = file.WriteString(`{"uuid":"uuid1","long_url":"http://example1.com","short_code":"code1"}
{"uuid":"uuid2","long_url":"http://example2.com","short_code":"code2"}`)
				require.NoError(t, err)
			},
			expected: result{
				memento: &Memento{
					State: []URL{
						{
							UUID:      "uuid1",
							LongURL:   "http://example1.com",
							ShortCode: "code1",
						},
						{
							UUID:      "uuid2",
							LongURL:   "http://example2.com",
							ShortCode: "code2",
						},
					},
				},
				err: nil,
			},
		},
		{
			name: "Invalid JSON",
			before: func(filePath string) {
				file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
				require.NoError(t, err)
				defer file.Close()

				_, err = file.WriteString(`{invalid json}`)
				require.NoError(t, err)
			},
			expected: result{
				memento: nil,
				err:     errors.ErrorFailedToReadFromFile,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := t.TempDir() + "/store-test.json"

			tt.before(filePath)

			fileRepo := NewFileRepository(filePath)
			memento, err := fileRepo.Load()

			if tt.expected.err != nil {
				assert.Equal(t, tt.expected.err, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.expected.memento, memento)
			}

			t.Cleanup(func() {
				os.RemoveAll(filePath)
			})
		})
	}
}

func Test_FileStorageRepository_Save(t *testing.T) {
	tests := []struct {
		name     string
		before   func(filePath string)
		payload  *Memento
		expected error
	}{
		{
			name:   "Success",
			before: func(_ string) {},
			payload: &Memento{
				State: []URL{
					{
						UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
						LongURL:   "http://example.com",
						ShortCode: "abcd1234",
					},
				},
			},
			expected: nil,
		},
		{
			name: "File not exists",
			before: func(filePath string) {
				os.Remove(filePath)
			},
			payload: &Memento{
				State: []URL{
					{
						UUID:      "6455bd07-e431-4851-af3c-4f703f726639",
						LongURL:   "http://example.com",
						ShortCode: "abcd1234",
					},
				},
			},
			expected: nil,
		},
		{
			name:   "Empty payload",
			before: func(_ string) {},
			payload: &Memento{
				State: []URL{},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := t.TempDir() + "/store-test.json"

			fileRepo := NewFileRepository(filePath)
			err := fileRepo.Save(tt.payload)

			assert.Equal(t, tt.expected, err)
		})
	}
}
