package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/xedom/codeduel/config"
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

func CreateCorsMiddleware(config *config.Config) Middleware {
	return func(next http.Handler) http.Handler {
		return CorsMiddleware(next, config)
	}
}

func CorsMiddleware(next http.Handler, config *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", config.CorsOrigin)
		w.Header().Set("Access-Control-Allow-Methods", config.CorsMethods)
		w.Header().Set("Access-Control-Allow-Headers", config.CorsHeaders)
		w.Header().Set("Access-Control-Allow-Credentials", fmt.Sprintf("%t", config.CorsCredentials))
		w.Header().Set("Access-Control-Expose-Headers", "Authorization")
		next.ServeHTTP(w, r)
	})
}

func OnlyInternalServiceMiddleware(config *config.Config, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-service-token")
		if token != config.ServiceToken {
			_ = WriteJSON(w, http.StatusUnauthorized, Error{Err: "Unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

type contextKey string

const AuthUser contextKey = "middleware.auth.user"

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("x-jwt-token")
		if tokenString == "" {
			// get from cookie
			cookie, err := r.Cookie("jwt")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				_ = WriteJSON(w, http.StatusUnauthorized, Error{Err: err.Error()})
				return
			}
			tokenString = cookie.Value
		}

		userHeader, err := utils.ValidateUserJWT(tokenString)
		if err != nil {
			_ = WriteJSON(w, http.StatusUnauthorized, Error{Err: err.Error()})
			return
		}

		ctx := context.WithValue(r.Context(), AuthUser, userHeader)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	}
}

func GetAuthUser(r *http.Request) *types.UserRequestHeader {
	return r.Context().Value(AuthUser).(*types.UserRequestHeader)
}
