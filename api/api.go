package api

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

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
	return &Server{
		config:  config,
		db:      db,
		address: fmt.Sprintf("%s:%s", config.Host, config.Port),
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
	v1.Handle("POST /validateToken", makeHTTPHandleFunc(s.handleValidateToken))
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
	// main.HandleFunc("/docs/", httpSwagger.Handler(httpSwagger.URL("http://"+s.address+"/docs/doc.json")))
	main.HandleFunc("/docs/", httpSwagger.Handler())
	main.Handle("/v1/", http.StripPrefix("/v1", v1))

	var wg sync.WaitGroup
	wg.Add(2)

	// HTTPS server
	go startHttpsServer(
		s.config,
		ChainMiddleware(CreateCorsMiddleware(s.config), LoggingMiddleware)(main),
		&wg,
	)

	// HTTP server
	go startHttpServer(
		s.config,
		ChainMiddleware(CreateCorsMiddleware(s.config), LoggingMiddleware)(main),
		&wg,
	)

	wg.Wait()

	return nil
}

func startHttpsServer(config *config.Config, handler http.Handler, wg *sync.WaitGroup) {
	tlsCert, err := tls.LoadX509KeyPair(config.SSLCert, config.SSLKey)
	if err != nil {
		log.Printf("%s%s failed to load SSL certificate: %s", utils.GetLogTag("API"), utils.GetLogTag("error"), err.Error())
	}

	sslCertFile := utils.GetEnv("SSL_CERT_FILE", "/etc/ssl/certs")
	log.Printf("%s SSL Cert file: %s", utils.GetLogTag("info"), sslCertFile)

	// certManager := autocert.Manager{
	// 	Prompt:     autocert.AcceptTOS,
	// 	HostPolicy: autocert.HostWhitelist("api.codeduel.it"),
	// 	Cache:      autocert.DirCache(sslCertFile),
	// 	// Cache:      autocert.DirCache("/etc/ssl/certs"),
	// }
	httpsAddress := fmt.Sprintf("%s:%s", config.Host, config.Port)

	server := &http.Server{
		Addr: httpsAddress,

		// setting timeouts to avoid Slowloris attack
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,

		Handler: handler,
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{tlsCert},
			// GetCertificate: certManager.GetCertificate,
			// MinVersion:     tls.VersionTLS12,
		},
	}

	log.Printf("%s server started on %s", utils.GetLogTag("HTTPS"), httpsAddress)
	err = server.ListenAndServeTLS("", "")

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("%s failed to start server: %s", utils.GetLogTag("HTTPS"), err.Error())
	} else if errors.Is(err, http.ErrServerClosed) {
		log.Printf("%s server closed", utils.GetLogTag("HTTPS"))
	} else {
		log.Printf("%s server started", utils.GetLogTag("HTTPS"))
	}

	wg.Done()
}

func startHttpServer(config *config.Config, handler http.Handler, wg *sync.WaitGroup) {
	httpAddress := fmt.Sprintf("%s:%s", config.Host, config.PortHttp)

	server := &http.Server{
		Addr:    httpAddress,
		Handler: handler,
		// Handler: certManager.HTTPHandler(nil),
	}

	log.Printf("%s server started on %s", utils.GetLogTag("HTTP"), httpAddress)
	err := server.ListenAndServe()

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("%s failed to start server: %s", utils.GetLogTag("HTTP"), err.Error())
	} else if errors.Is(err, http.ErrServerClosed) {
		log.Printf("%s server closed", utils.GetLogTag("HTTP"))
	} else {
		log.Printf("%s server started", utils.GetLogTag("HTTP"))
	}

	wg.Done()
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
