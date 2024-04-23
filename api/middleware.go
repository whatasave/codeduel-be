package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

type Middleware func(http.Handler) http.Handler

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) WriteHeader(status int) {
	w.statusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedWriter{w, http.StatusOK}
		next.ServeHTTP(wrapped, r)

		log.Printf("%s %d %s %s %v",
			utils.GetLogTag("api"),
			wrapped.statusCode,
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

func ChainMiddleware(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, x-jwt-token")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

type contextKey string
const AuthUser contextKey = "middleware.auth.user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("x-jwt-token")
		if tokenString == "" {
			// get from cookie
			cookie, err := r.Cookie("jwt")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				WriteJSON(w, http.StatusUnauthorized, Error{Err: err.Error()})
				return
			}
			tokenString = cookie.Value
		}

		userHeader, err := utils.ValidateUserJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, Error{Err: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), AuthUser, userHeader)
		r = r.WithContext(ctx)

		r.Header.Set("x-user-id", fmt.Sprintf("%d", userHeader.ID))
		r.Header.Set("x-user-username", userHeader.Username)
		r.Header.Set("x-user-email", userHeader.Email)

		next.ServeHTTP(w, r)
	})
}

func GetAuthUser(r *http.Request) *types.UserRequestHeader {
	return r.Context().Value(AuthUser).(*types.UserRequestHeader)
}
