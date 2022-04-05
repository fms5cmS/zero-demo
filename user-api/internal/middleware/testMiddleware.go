package middleware

import (
	"fmt"
	"net/http"
)

type TestMiddleware struct {
}

func NewTestMiddleware() *TestMiddleware {
	return &TestMiddleware{}
}

func (m *TestMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("middleware start")
		// Passthrough to next handler if need
		next(w, r)
		fmt.Println("middleware end")
	}
}
