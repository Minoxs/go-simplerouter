package simplerouter

import (
	"net/http"
)

// SimpleMiddleware is the function signature that is supposed to be used with this library's methods.
// It allows for both-way handling.
type SimpleMiddleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// middlewareRouter is the router returned when a middleware is registered.
// This struct is internal and should not be created manually, do it at your own risk.
// Registered middlewares are added to a write-only stack, and once a method is registered
// the methods are peeked from the stack one-by-one and the middleware chain is created.
type middlewareRouter struct {
	*SimpleRouter
	middlewareStack []SimpleMiddleware
}

// newMiddlewareRouter creates a middleware router with default settings
func newMiddlewareRouter(root *SimpleRouter) *middlewareRouter {
	return &middlewareRouter{
		SimpleRouter:    root,
		middlewareStack: make([]SimpleMiddleware, 0),
	}
}

// push adds a new middleware to the stack
func (r *middlewareRouter) push(m SimpleMiddleware) *middlewareRouter {
	r.middlewareStack = append(r.middlewareStack, m)
	return r
}

// createChain creates the middleware chain, returning a handler function
func (r *middlewareRouter) createChain(h http.HandlerFunc) http.HandlerFunc {
	chain := h

	for i := len(r.middlewareStack) - 1; i >= 0; i-- {
		// This looks really complicated
		// But it's just a janky way of moving the thingies into another scope
		// So that you don't get a middleware referencing itself or something silly like that
		chain = func(current http.HandlerFunc, top SimpleMiddleware) http.HandlerFunc {
			return func(wrt http.ResponseWriter, req *http.Request) {
				top(wrt, req, current)
			}
		}(chain, r.middlewareStack[i])
	}

	return chain
}

// FileServer creates a file server with files from root.
// Be careful with this thing. Check out http.FileServer for
// more information.
func (s *middlewareRouter) FileServer(root string) *middlewareRouter {
	s.fileServer = s.createChain(http.FileServer(http.Dir(root)).ServeHTTP)
	return s
}

// Get registers a get endpoint
func (s *middlewareRouter) Get(pattern string, handle http.HandlerFunc) *middlewareRouter {
	f := s.createChain(handle)
	s.SimpleRouter.Get(pattern, f)
	return s
}

// Post registers a post endpoint
func (s *middlewareRouter) Post(pattern string, handle http.HandlerFunc) *middlewareRouter {
	f := s.createChain(handle)
	s.SimpleRouter.Post(pattern, f)
	return s
}

// Delete register a delete endpoint
func (s *middlewareRouter) Delete(pattern string, handle http.HandlerFunc) *middlewareRouter {
	f := s.createChain(handle)
	s.SimpleRouter.Delete(pattern, f)
	return s
}

// Use registers a middleware is run in the order they are registered.
// After every middleware the handler is run
func (s *middlewareRouter) Use(m SimpleMiddleware) *middlewareRouter {
	s.push(m)
	return s
}
