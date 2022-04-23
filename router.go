package simplerouter

import (
	"fmt"
	"net/http"
)

// SimpleRouter is a HTTP Handler, and is used to simplify creating endpoints.
// While there is no problem with manually creating the struct, using New is
// highly encouraged.
type SimpleRouter struct {
	prefix     string
	mux        *SimpleMux
	fileServer http.Handler
}

// New creates a new SimpleRouter with default settings
func New() *SimpleRouter {
	return &SimpleRouter{
		prefix:     "",
		mux:        NewSimpleMux(),
		fileServer: nil,
	}
}

// ServeHTTP is implemented so that SimpleRouter implements http.Handler
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

// FileServer creates a file server with files from root.
// Be careful with this thing. Check out http.FileServer for
// more information.
func (s *SimpleRouter) FileServer(root string) *SimpleRouter {
	s.fileServer = http.FileServer(http.Dir(root))
	return s
}

// Prefix creates a subrouter that where every endpoint stemming from it
// has to start with... you guessed it! the prefix that was set.
func (s *SimpleRouter) Prefix(prefix string) *SimpleRouter {
	return &SimpleRouter{
		prefix:     s.prefix + cleanPattern(prefix),
		mux:        s.mux,
		fileServer: s.fileServer,
	}
}

// Get registers a get endpoint
func (s *SimpleRouter) Get(pattern string, handle http.HandlerFunc) *SimpleRouter {
	pattern = cleanPattern(pattern)
	s.mux.HandleFunc(GET, s.prefix+pattern, handle)
	return s
}

// Post registers a post endpoint
func (s *SimpleRouter) Post(pattern string, handle http.HandlerFunc) *SimpleRouter {
	pattern = cleanPattern(pattern)
	s.mux.HandleFunc(POST, s.prefix+pattern, handle)
	return s
}

// Delete register a delete endpoint
func (s *SimpleRouter) Delete(pattern string, handle http.HandlerFunc) *SimpleRouter {
	pattern = cleanPattern(pattern)
	s.mux.HandleFunc(DELETE, s.prefix+pattern, handle)
	return s
}

// Use registers a middleware is run in the order they are registered.
// After every middleware the handler is run
func (s *SimpleRouter) Use(m SimpleMiddleware) *middlewareRouter {
	return newMiddlewareRouter(s).push(m)
}
