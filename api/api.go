package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/xedom/codeduel/db"
	"github.com/xedom/codeduel/utils"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	host       string
	port       string
	listenAddr string
	db         db.DB
}

type ApiError struct {
	Err string `json:"error"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func NewAPIServer(host, port string, db db.DB) *APIServer {
	address := fmt.Sprintf("%s:%s", host, port)
	log.Print("[API] Starting API server on ", address)
	return &APIServer{
		host:       host,
		port:       port,
		listenAddr: address,
		db:         db,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1", makeHTTPHandleFunc(s.handleRoot))
	router.HandleFunc("/api/v1/health", makeHTTPHandleFunc(s.handleHealth))
	router.HandleFunc("/api/v1/user", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc("/api/v1/user/{id}", makeHTTPHandleFunc(s.handleUserByID))
	router.HandleFunc("/api/v1/profile", authMiddleware(makeHTTPHandleFunc(s.handleProfile)))

	router.HandleFunc("/api/v1/validateToken", makeHTTPHandleFunc(s.handleValidateToken)) // TODO: make it accessible only by lobby service
	// router.HandleFunc("/api/v1/lobbies", makeHTTPHandleFunc(s.handleLobbies))

	router.HandleFunc("/api/v1/auth/github", makeHTTPHandleFunc(s.handleGithubAuth))
	router.HandleFunc("/api/v1/auth/github/callback", makeHTTPHandleFunc(s.handleGithubAuthCallback))

	frontendUrl := os.Getenv("FRONTEND_URL")

	err := http.ListenAndServe(s.listenAddr, handlers.CORS(
		handlers.AllowedOrigins([]string{frontendUrl}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With", "x-jwt-token"}),
		handlers.AllowCredentials(),
	)(router))

	if err != nil {
		log.Fatal("[API] Cannot start http server: ", err)
	}
}

func (s *APIServer) handleRoot(w http.ResponseWriter, r *http.Request) error {
	// host := fmt.Sprintf("http://%s:%s", s.host, s.port)
	host := fmt.Sprintf("http://%s", r.Host)

	return WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Welcome to CodeDuel API",
		"version": "v1",
		"status":  "ok",
		"apis": []string{
			fmt.Sprintf("%s/api/v1", host),
			fmt.Sprintf("%s/api/v1/health", host),
			fmt.Sprintf("%s/api/v1/user", host),
			fmt.Sprintf("%s/api/v1/user/{id}", host),
			fmt.Sprintf("%s/api/v1/profile", host),
			// fmt.Sprintf("%s/api/v1/lobbies", host),
			fmt.Sprintf("%s/api/v1/auth/github", host),
			fmt.Sprintf("%s/api/v1/auth/github/callback", host),
		},
	})
}

func (s *APIServer) handleHealth(w http.ResponseWriter, r *http.Request) error {
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func makeHTTPHandleFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			WriteJSON(w, http.StatusInternalServerError, ApiError{Err: err.Error()})
		}
	}
}

func authMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("x-jwt-token")
		if tokenString == "" {
			// get from cookie
			cookie, err := r.Cookie("jwt")
			if err != nil {
				WriteJSON(w, http.StatusUnauthorized, ApiError{Err: err.Error()})
				return
			}
			tokenString = cookie.Value
		}

		userHeader, err := utils.ValidateUserJWT(tokenString)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Err: err.Error()})
			return
		}

		r.Header.Set("x-user-id", fmt.Sprintf("%d", userHeader.ID))
		r.Header.Set("x-user-username", userHeader.Username)
		r.Header.Set("x-user-email", userHeader.Email)

		handlerFunc(w, r)
	}
}
