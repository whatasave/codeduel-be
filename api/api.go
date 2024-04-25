package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/xedom/codeduel/config"
	"github.com/xedom/codeduel/db"
	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/xedom/codeduel/docs"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type Server struct {
	config  *config.Config
	address string
	db      db.DB
}

type Error struct {
	Err string `json:"error"`
}

type Handler func(w http.ResponseWriter, r *http.Request) error

func NewAPIServer(config *config.Config, db db.DB) *Server {
	address := fmt.Sprintf("%s:%s", config.Host, config.Port)
	log.Printf("%s Starting API server on http://%s", utils.GetLogTag("main"), address)
	log.Printf("%s Docs http://%s/docs", utils.GetLogTag("main"), address)
	return &Server{
		config:  config,
		db:      db,
		address: address,
	}
}

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

// @host		localhost
// @schemes	http
func (s *Server) Run() error {
	v1 := http.NewServeMux()
	v1.Handle("/user", s.GetUserRouter())
	v1.Handle("/user/", s.GetUserRouter())
	v1.Handle("/lobby", s.GetLobbyRouter())
	v1.Handle("/lobby/", s.GetLobbyRouter())
	v1.Handle("/challenge", s.GetChallengeRouter())
	v1.Handle("/challenge/", s.GetChallengeRouter())
	v1.Handle("/auth/github", s.GetGithubAuthRouter())
	v1.Handle("/auth/github/", s.GetGithubAuthRouter())

	main := http.NewServeMux()
	main.HandleFunc("/v1", makeHTTPHandleFunc(s.handleRoot))
	main.HandleFunc("/health", makeHTTPHandleFunc(s.handleHealth))
	main.HandleFunc("POST /validateToken", makeHTTPHandleFunc(s.handleValidateToken)) // TODO: make it accessible only by lobby service
	// main.HandleFunc("/docs/", httpSwagger.Handler(httpSwagger.URL("http://"+s.address+"/docs/doc.json")))
	main.HandleFunc("/docs/", httpSwagger.Handler())
	main.Handle("/v1/", http.StripPrefix("/v1", v1))
	
	serverSSL := &http.Server{
		Addr:    s.address,
		Handler: ChainMiddleware(CorsMiddleware, LoggingMiddleware)(main),
		TLSConfig: &tls.Config{},
	}

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", s.config.Host, s.config.PortHttp),
		Handler: ChainMiddleware(CorsMiddleware, LoggingMiddleware)(main),
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		log.Printf("%s HTTPS server starting...", utils.GetLogTag("info"))
		err := serverSSL.ListenAndServeTLS(s.config.SSLCert, s.config.SSLKey)

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("%s failed to start HTTPS server: %s", utils.GetLogTag("error"), err.Error())
		} else if errors.Is(err, http.ErrServerClosed) {
			log.Printf("%s HTTPS server closed", utils.GetLogTag("error"))
		} else {
			log.Printf("%s HTTPS server started", utils.GetLogTag("info"))
		}

		wg.Done()
	}()

	go func() {
		log.Printf("%s HTTP server starting...", utils.GetLogTag("info"))
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("%s failed to start HTTP server: %s", utils.GetLogTag("error"), err.Error())
		} else if errors.Is(err, http.ErrServerClosed) {
			log.Printf("%s HTTP server closed", utils.GetLogTag("error"))
		} else {
			log.Printf("%s HTTP server started", utils.GetLogTag("info"))
		}

		wg.Done()
	}()

	wg.Wait()

	return nil
}


// @Summary		Root
// @Description	Root endpoint
// @Tags			root
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]any
// @Router			/v1 [get]
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) error {
	host := fmt.Sprintf("http://%s", r.Host)
	swaggerUrl := fmt.Sprintf("%s/docs/index.html", host)

	return WriteJSON(w, http.StatusOK, map[string]any{
		"message": "Welcome to CodeDuel API",
		"version": "v1",
		"status":  "ok",
		"apis":    swaggerUrl,
	})
}

// @Summary		Health check
// @Description	Health check endpoint
// @Tags			root
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]string
// @Router			/health [get]
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) error {
	// return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	return WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// @Summary		Validate JWT Token
// @Description	Validate if the user JWT token is valid, and return user data. Used from other services to validate user token
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			token	body		types.VerifyToken	true	"Service token"
// @Success		200		{object}	types.User
// @Failure		500		{object}	Error
// @Router			/validateToken [post]
func (s *Server) handleValidateToken(w http.ResponseWriter, r *http.Request) error {
	verifyTokenBody := &types.VerifyToken{}
	if err := json.NewDecoder(r.Body).Decode(verifyTokenBody); err != nil {
		return err
	}

	decodedUserData, err := utils.ValidateUserJWT(verifyTokenBody.JWTToken)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	// user, err := s.db.GetUserByID(userID)
	// if err != nil {
	// 	return err
	// }

	return WriteJSON(w, http.StatusOK, decodedUserData)
}

func makeHTTPHandleFunc(fn Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %s", utils.GetLogTag("error"), err.Error())
			err := WriteJSON(w, http.StatusInternalServerError, Error{Err: err.Error()})
			if err != nil {
				log.Printf("%s %s", utils.GetLogTag("error"), err.Error())
			}
		}
	}
}
