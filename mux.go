package simplerouter

import "net/http"

type (
	method   uint8
	handlers [3]http.Handler
)

const (
	GET    method = 0
	POST          = 1
	DELETE        = 2
)

type SimpleMux struct {
	m map[string]*handlers
}

func NewSimpleMux() *SimpleMux {
	return &SimpleMux{
		m: make(map[string]*handlers),
	}
}

func (s *SimpleMux) HandleFunc(meth method, url string, handler http.Handler) {
	_, ok := s.m[url]
	if !ok {
		s.m[url] = &handlers{nil, nil, nil}
	}
	s.m[url][meth] = handler
}

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
