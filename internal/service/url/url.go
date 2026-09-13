package url

import (
	"errors"
	"fmt"

	"github.com/annpolukaro/url-shortener/internal/lib/random"
	"github.com/annpolukaro/url-shortener/internal/storage"
)

const aliasLength = 6

type URLStorage interface {
	SaveURL(urlSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
}

type Service struct {
	storage URLStorage
}

func New(storage URLStorage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) SaveURL(urlToSave string, alias string) (string, error) {
	if alias != "" {
		_, err := s.storage.SaveURL(urlToSave, alias)
		if err != nil {
			return "", fmt.Errorf("save url: %w", err)
		}
		return alias, nil
	}
	for {
		alias = random.NewRandomString(aliasLength)

		_, err := s.storage.SaveURL(urlToSave, alias)

		if errors.Is(err, storage.ErrURLExists) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("save url: %w", err)
		}
		return alias, nil
	}
}

func (s *Service) GetURL(alias string) (string, error) {
	url, err := s.storage.GetURL(alias)
	if err != nil {
		return "", fmt.Errorf("get url: %w", err)

	}
	return url, nil
}
