// Package geocode provides a client for the Google Geocoding API
package geocode

import "net/http"

const baseURL = "https://maps.googleapis.com/maps/api/geocode"

type Service struct {
	client *http.Client
	key    string
	url    string
}

// NewService creates a new geocode service with the given http client and Google Geocoding API key
func NewService(client *http.Client, key string) *Service {
	return &Service{
		client: client,
		key:    key,
		url:    baseURL,
	}
}

// SetURL allows overwriting the base url
func (s *Service) SetURL(url string) {
	s.url = url
}
