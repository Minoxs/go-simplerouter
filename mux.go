package simplerouter

import "net/http"

// Helper types
type (
	method   uint8
	handlers [3]http.Handler
)

// Helper constants for accessing handler array
const (
	GET    method = 0
	POST          = 1
	DELETE        = 2
)

// SimpleMux is a simple map-based mux.
// Paths are mapped to an array with each possible method's handler.
// Use NewSimpleMux to create.
type SimpleMux struct {
	m map[string]*handlers
}

// NewSimpleMux creates a SimpleMux with default settings (there are no settings)
func NewSimpleMux() *SimpleMux {
	return &SimpleMux{
		m: make(map[string]*handlers),
	}
}

// HandleFunc register a function handler for a specific method
func (s *SimpleMux) HandleFunc(meth method, url string, handler http.Handler) {
	_, ok := s.m[url]
	if !ok {
		s.m[url] = &handlers{nil, nil, nil}
	}
	s.m[url][meth] = handler
}

// GetHandler returns the handler method for a given request.
// Also returns if the method is allowed or not.
// If the handler method is nil, then there is no endpoint registered.
func (s *SimpleMux) GetHandler(req *http.Request) (http.Handler, bool) {
	h, ok := s.m[req.URL.Path]
	if !ok {
		return nil, true
	}

	switch req.Method {
	case http.MethodGet:
		return h[GET], h[GET] != nil
	case http.MethodPost:
		return h[POST], h[POST] != nil
	case http.MethodDelete:
		return h[DELETE], h[DELETE] != nil
	default:
		return nil, false
	}
}
