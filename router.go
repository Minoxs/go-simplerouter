package simplerouter

import (
	"fmt"
	"net/http"
)

type SimpleRouter struct {
	prefix     string
	mux        *SimpleMux
	fileServer http.Handler
}

func New() *SimpleRouter {
	return &SimpleRouter{
		prefix:     "",
		mux:        NewSimpleMux(),
		fileServer: nil,
	}
}

func (s *SimpleRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Try handling using mux
	handler, allowed := s.mux.GetHandler(r)
	if !allowed {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = fmt.Fprintf(w, "%d Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Handle request
	if handler != nil {
		handler.ServeHTTP(w, r)
		return
	}

	// If there is a file server assigned, try it out
	if s.fileServer != nil {
		s.fileServer.ServeHTTP(w, r)
		return
	}

	// Only gets here if path not found
	http.NotFound(w, r)
}

func (s *SimpleRouter) FileServer(root string) *SimpleRouter {
	s.fileServer = http.FileServer(http.Dir(root))
	return s
}

func (s *SimpleRouter) Prefix(prefix string) *SimpleRouter {
	return &SimpleRouter{
		prefix:     s.prefix + cleanPattern(prefix),
		mux:        s.mux,
		fileServer: s.fileServer,
	}
}

func (s *SimpleRouter) Get(pattern string, handle http.HandlerFunc) *SimpleRouter {
	pattern = cleanPattern(pattern)
	s.mux.HandleFunc(GET, s.prefix+pattern, handle)
	return s
}

func (s *SimpleRouter) Post(pattern string, handle http.HandlerFunc) *SimpleRouter {
	pattern = cleanPattern(pattern)
	s.mux.HandleFunc(POST, s.prefix+pattern, handle)
	return s
}

func (s *SimpleRouter) Delete(pattern string, handle http.HandlerFunc) *SimpleRouter {
	pattern = cleanPattern(pattern)
	s.mux.HandleFunc(DELETE, s.prefix+pattern, handle)
	return s
}

func (s *SimpleRouter) Use(m SimpleMiddleware) *middlewareRouter {
	return newMiddlewareRouter(s).push(m)
}
