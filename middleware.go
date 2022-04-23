package simplerouter

import (
	"net/http"
)

type SimpleMiddleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

type middlewareRouter struct {
	*SimpleRouter
	middlewareStack []SimpleMiddleware
}

func newMiddlewareRouter(root *SimpleRouter) *middlewareRouter {
	return &middlewareRouter{
		SimpleRouter:    root,
		middlewareStack: make([]SimpleMiddleware, 0),
	}
}

func (r *middlewareRouter) push(m SimpleMiddleware) *middlewareRouter {
	r.middlewareStack = append(r.middlewareStack, m)
	return r
}

func (r *middlewareRouter) createChain(h http.HandlerFunc) http.HandlerFunc {
	chain := h

	for i := len(r.middlewareStack) - 1; i >= 0; i-- {
		chain = func(current http.HandlerFunc, top SimpleMiddleware) http.HandlerFunc {
			return func(wrt http.ResponseWriter, req *http.Request) {
				top(wrt, req, current)
			}
		}(chain, r.middlewareStack[i])
	}

	return chain
}

func (s *middlewareRouter) FileServer(root string) *middlewareRouter {
	s.fileServer = s.createChain(http.FileServer(http.Dir(root)).ServeHTTP)
	return s
}

func (s *middlewareRouter) Get(pattern string, handle http.HandlerFunc) *middlewareRouter {
	f := s.createChain(handle)
	s.SimpleRouter.Get(pattern, f)
	return s
}

func (s *middlewareRouter) Post(pattern string, handle http.HandlerFunc) *middlewareRouter {
	f := s.createChain(handle)
	s.SimpleRouter.Post(pattern, f)
	return s
}

func (s *middlewareRouter) Delete(pattern string, handle http.HandlerFunc) *middlewareRouter {
	f := s.createChain(handle)
	s.SimpleRouter.Delete(pattern, f)
	return s
}

func (s *middlewareRouter) Use(m SimpleMiddleware) *middlewareRouter {
	s.push(m)
	return s
}
