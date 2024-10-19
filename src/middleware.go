package src

import (
	"net/http"
)

type Middleware func(next http.Handler) http.Handler

func RegisterMiddlewares(mx http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		mx = middleware(mx)
	}

	return mx
}

func GetUserFromRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != "" {
			user, err := ParseToken(token)
			if err != nil {
				ErrorHandler(w, err)
				return
			}
			r = SetRequestContext(r, UserContextKey, user)
		}
		next.ServeHTTP(w, r)
	})
}
