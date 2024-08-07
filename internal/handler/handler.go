package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rabboni171/url-shortener/configs"
	"github.com/rabboni171/url-shortener/internal/service"
	"github.com/rs/zerolog"
)

type URLHandler struct {
	service service.IURLService
	log     *zerolog.Logger
	cfg     *configs.Config
}

func NewURLHandler(cfg *configs.Config, log *zerolog.Logger, service service.IURLService) *URLHandler {
	return &URLHandler{
		log:     log,
		service: service,
		cfg:     cfg,
	}
}

func (h *URLHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var request struct {
		URL string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.Shorten(request.URL)
	if err != nil {
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	str := []string{h.cfg.AppParams.BaseURL, h.cfg.AppParams.PortRun}
	baseURL := strings.Join(str, ":")
	fullURL := fmt.Sprintf("%s/%s", baseURL, shortURL)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"shortURL": fullURL})
}

func (h *URLHandler) ResolveURL(w http.ResponseWriter, r *http.Request) {
	shortURL := mux.Vars(r)["shortURL"]

	originalURL, err := h.service.Resolve(shortURL)
	if err != nil {
		http.Error(w, "Failed to resolve URL", http.StatusInternalServerError)
		return
	}
	if originalURL == "" {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
