package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/xedom/codeduel/config"
	"github.com/xedom/codeduel/db"
	"github.com/xedom/codeduel/utils"

	_ "github.com/swaggo/http-swagger/example/gorilla/docs"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type APIServer struct {
	config    *config.Config
	listenAddr string
	db         db.DB
}

type ApiError struct {
	Err string `json:"error"`
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func NewAPIServer(config *config.Config, db db.DB) *APIServer {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Print("[API] Starting API server on http://", address)
	log.Print("[API] Docs http://localhost:5000/docs/index.html")
	return &APIServer{
		config:     config,
		db:         db,
		listenAddr: address,
	}
}

//	@title						CodeDuel API
//	@version					1.0
//	@description				Backend API for CodeDuel
//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						token
//	@termsOfService				http://swagger.io/terms/
//	@contact.name				API Support
//	@contact.email				support@codeduel
//	@license.name				Apache 2.0
//	@license.url				http://www.apache.org/licenses/LICENSE-2.0.html
//	@host						localhost:5000
//	@schemes					http 
//	@BasePath					/v1/
func (s *APIServer) Run() {
	router := mux.NewRouter()
	
	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)

	router.HandleFunc("/v1", makeHTTPHandleFunc(s.handleRoot))

	router.HandleFunc("/v1/health", makeHTTPHandleFunc(s.handleHealth))
	router.HandleFunc("/v1/validateToken", makeHTTPHandleFunc(s.handleValidateToken)) // TODO: make it accessible only by lobby service

	router.HandleFunc("/v1/user", makeHTTPHandleFunc(s.handleUser))
	router.HandleFunc("/v1/user/{id}", makeHTTPHandleFunc(s.handleUserByID))
	router.HandleFunc("/v1/profile", authMiddleware(makeHTTPHandleFunc(s.handleProfile)))

	router.HandleFunc("/v1/challenge", makeHTTPHandleFunc(s.handleGetChallenges)).Methods(http.MethodGet)
	router.HandleFunc("/v1/challenge", authMiddleware(makeHTTPHandleFunc(s.handleCreateChallenge))).Methods(http.MethodPost)
	router.HandleFunc("/v1/challenge/{id}", makeHTTPHandleFunc(s.handleGetChallengeByID)).Methods(http.MethodGet)	
	router.HandleFunc("/v1/challenge/{id}", authMiddleware(makeHTTPHandleFunc(s.handleUpdateChallenge))).Methods(http.MethodPut)
	router.HandleFunc("/v1/challenge/{id}", authMiddleware(makeHTTPHandleFunc(s.handleDeleteChallenge))).Methods(http.MethodDelete)

	router.HandleFunc("/v1/auth/github", makeHTTPHandleFunc(s.handleGithubAuth))
	router.HandleFunc("/v1/auth/github/callback", makeHTTPHandleFunc(s.handleGithubAuthCallback))

	err := http.ListenAndServe(s.listenAddr, handlers.CORS(
		handlers.AllowedOrigins([]string{s.config.FrontendURL}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Access-Control-Allow-Headers", "Authorization", "X-Requested-With", "x-jwt-token"}),
		handlers.AllowCredentials(),
	)(router))

	if err != nil {
		log.Fatal("[API] Cannot start http server: ", err)
	}
}

func (s *APIServer) handleRoot(w http.ResponseWriter, r *http.Request) error {
	host := fmt.Sprintf("http://%s", r.Host)

	return WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Welcome to CodeDuel API",
		"version": "v1",
		"status":  "ok",
		"apis": []string{
			fmt.Sprintf("%s/v1", host),
			fmt.Sprintf("%s/v1/health", host),
			fmt.Sprintf("%s/v1/user", host),
			fmt.Sprintf("%s/v1/user/{id}", host),
			fmt.Sprintf("%s/v1/profile", host),
			// fmt.Sprintf("%s/v1/lobbies", host),
			fmt.Sprintf("%s/v1/auth/github", host),
			fmt.Sprintf("%s/v1/auth/github/callback", host),
		},
	})
}

// health check
//	@Summary		Health check
//	@Description	Health check
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/v1/health [get]
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
