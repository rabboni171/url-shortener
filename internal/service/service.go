package service

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"

	"github.com/rabboni171/url-shortener/internal/repository"
)

type URLServiceImpl struct {
	repo repository.IURLRepository
}

func NewURLServiceImpl(repo repository.IURLRepository) IURLService {
	return &URLServiceImpl{
		repo: repo,
	}
}

func (s *URLServiceImpl) Shorten(url string) (string, error) {
	hash := sha256.Sum256([]byte(url))
	shortURL := base64.URLEncoding.EncodeToString(hash[:])[:8]

	err := s.repo.Save(shortURL, url)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *URLServiceImpl) Resolve(shortURL string) (string, error) {
	originalURL, err := s.repo.Get(shortURL)
	if err != nil {
		return "", err
	}

	if originalURL == "" {
		return "", errors.New("URL not found")
	}

	return originalURL, nil
}
