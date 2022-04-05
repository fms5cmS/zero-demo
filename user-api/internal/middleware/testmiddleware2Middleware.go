package middleware

import "net/http"

type TestMiddleware2Middleware struct {
}

func NewTestMiddleware2Middleware() *TestMiddleware2Middleware {
	return &TestMiddleware2Middleware{}
}

func (m *TestMiddleware2Middleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation

		// Passthrough to next handler if need
		next(w, r)
	}
}
