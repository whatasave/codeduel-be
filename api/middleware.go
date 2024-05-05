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

type Middleware2 func(w http.ResponseWriter, r *http.Request) *http.Request

func AuthMiddleware(w http.ResponseWriter, r *http.Request) *http.Request {
	tokenString := r.Header.Get("x-token")
	if tokenString == "" {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			_ = WriteJSON(w, http.StatusUnauthorized, Error{Err: err.Error()})
			return r
		}
		tokenString = cookie.Value
	}

	log.Printf("tokenString: %s", tokenString)

	userHeader, err := utils.ValidateUserJWT(tokenString)
	if err != nil {
		_ = WriteJSON(w, http.StatusUnauthorized, Error{Err: err.Error()})
		return r
	}

	ctx := context.WithValue(r.Context(), AuthUser, userHeader)
	r = r.WithContext(ctx)

	return r
}

func GetAuthUser(r *http.Request) *types.UserRequestHeader {
	user := r.Context().Value(AuthUser)
	if user == nil {
		log.Printf("User is nil")
		return nil
	}

	return user.(*types.UserRequestHeader)
}
