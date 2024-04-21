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
	"github.com/xedom/codeduel/docs"
	"github.com/xedom/codeduel/utils"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/xedom/codeduel/docs"
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
	log.Printf("%s Starting API server on http://%s", utils.GetLogTag("main"), address)
	log.Printf("%s Docs http://%s/docs/index.html", utils.GetLogTag("main"), address)
	return &APIServer{
		config:     config,
		db:         db,
		listenAddr: address,
	}
}

// https://github.com/swaggo/swag?tab=readme-ov-file#general-api-info

//	@title			CodeDuel API
//	@version		1.0
//	@description	Backend API for CodeDuel
//	@termsOfService	http://swagger.io/terms/

//	@securityDefinitions.basic	BasicAuth

//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						token
//	@description				Authorization token


//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@codeduel

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

//	@host		127.0.0.1:5000
//	@schemes	http
//	@basePath	/v1
func (s *APIServer) Run() {
	router := mux.NewRouter()

	docs.SwaggerInfo.Host = s.config.Host + ":" + s.config.Port
	httpSwagger.URL("http://" + s.config.Host + ":" + s.config.Port + "/docs/swagger.json")

	router.PathPrefix("/docs").Handler(httpSwagger.WrapHandler)

	router.HandleFunc("/v1", makeHTTPHandleFunc(s.handleRoot))

	router.HandleFunc("/v1/health", makeHTTPHandleFunc(s.handleHealth)).Methods(http.MethodGet)
	router.HandleFunc("/v1/validateToken", makeHTTPHandleFunc(s.handleValidateToken)).Methods(http.MethodPost) // TODO: make it accessible only by lobby service

	router.HandleFunc("/v1/user", makeHTTPHandleFunc(s.handleGetUsers)).Methods(http.MethodGet)
	router.HandleFunc("/v1/user", makeHTTPHandleFunc(s.handleCreateUser)).Methods(http.MethodPost)
	router.HandleFunc("/v1/user/{username}", makeHTTPHandleFunc(s.handleGetUserByUsername)).Methods(http.MethodGet)
	router.HandleFunc("/v1/user/{username}", authMiddleware(makeHTTPHandleFunc(s.handleDeleteUserByUsername))).Methods(http.MethodDelete)
	router.HandleFunc("/v1/profile", authMiddleware(makeHTTPHandleFunc(s.handleProfile))).Methods(http.MethodGet)

	router.HandleFunc("/v1/challenge", makeHTTPHandleFunc(s.handleGetChallenges)).Methods(http.MethodGet)
	router.HandleFunc("/v1/challenge/{id}", makeHTTPHandleFunc(s.handleGetChallengeByID)).Methods(http.MethodGet)	
	router.HandleFunc("/v1/challenge", authMiddleware(makeHTTPHandleFunc(s.handleCreateChallenge))).Methods(http.MethodPost)
	router.HandleFunc("/v1/challenge/{id}", authMiddleware(makeHTTPHandleFunc(s.handleUpdateChallenge))).Methods(http.MethodPut)
	router.HandleFunc("/v1/challenge/{id}", authMiddleware(makeHTTPHandleFunc(s.handleDeleteChallenge))).Methods(http.MethodDelete)

	router.HandleFunc("/v1/auth/github", makeHTTPHandleFunc(s.handleGithubAuth)).Methods(http.MethodGet)
	router.HandleFunc("/v1/auth/github/callback", makeHTTPHandleFunc(s.handleGithubAuthCallback)).Methods(http.MethodGet)

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

//	@Summary		Root
//	@Description	Root endpoint
//	@Tags			root
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]any
//	@Router			/v1 [get]
func (s *APIServer) handleRoot(w http.ResponseWriter, r *http.Request) error {
	host := fmt.Sprintf("http://%s", r.Host)
	swaggerUrl := fmt.Sprintf("%s/docs/index.html", host)

	return WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Welcome to CodeDuel API",
		"version": "v1",
		"status":  "ok",
		"apis": swaggerUrl,
	})
}

//	@Summary		Health check
//	@Description	Health check endpoint
//	@Tags			root
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
			log.Printf("%s%s Endpoint: %s Error: %s", utils.GetLogTag("api"), utils.GetLogTag("error"), r.RequestURI, err.Error())
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
