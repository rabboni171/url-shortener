package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rabboni171/url-shortener/configs"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockURLService struct {
	mock.Mock
}

func (m *MockURLService) Shorten(url string) (string, error) {
	args := m.Called(url)
	return args.String(0), args.Error(1)
}

func (m *MockURLService) Resolve(shortURL string) (string, error) {
	args := m.Called(shortURL)
	return args.String(0), args.Error(1)
}

func TestShortenURL(t *testing.T) {
	mockSvc := new(MockURLService)
	mockSvc.On("Shorten", "https://example.com").Return("shortURL", nil)

	// Создание конфигурации и логгера для тестов
	cfg := &configs.Config{
		AppParams: configs.Params{
			ServerURL: "http://localhost:8080",
		},
	}
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Logger()

	h := &URLHandler{
		service: mockSvc,
		log:     &logger,
		cfg:     cfg,
	}

	reqBody, _ := json.Marshal(map[string]string{
		"url": "https://example.com",
	})
	req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(reqBody))
	rec := httptest.NewRecorder()

	h.ShortenURL(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "http://localhost:8080/shortURL", resp["shortURL"])

	mockSvc.AssertExpectations(t)
}

func TestResolveURL(t *testing.T) {
	mockSvc := new(MockURLService)
	mockSvc.On("Resolve", "shortURL").Return("https://example.com", nil)

	// Создание конфигурации и логгера для тестов
	cfg := &configs.Config{
		AppParams: configs.Params{
			ServerURL: "http://localhost:8080",
		},
	}
	logger := zerolog.New(zerolog.NewConsoleWriter()).With().Logger()

	h := &URLHandler{
		service: mockSvc,
		log:     &logger,
		cfg:     cfg,
	}

	req := httptest.NewRequest(http.MethodGet, "/resolve?shortURL=shortURL", nil)
	rec := httptest.NewRecorder()

	h.ResolveURL(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp map[string]string
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "https://example.com", resp["originalURL"])

	mockSvc.AssertExpectations(t)
}
