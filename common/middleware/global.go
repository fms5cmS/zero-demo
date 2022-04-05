package middleware

import (
	"fmt"
	"net/http"
)

type GlobalMiddleware struct {
}

func NewGlobalMiddleware() *GlobalMiddleware {
	return &GlobalMiddleware{}
}

func (m *GlobalMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("全局中间件 start")
		next(w, r)
		fmt.Println("全局中间件 end")
	}
}
