package config

import (
	"github.com/xedom/codeduel/utils"
)

type Config struct {
	Host                            string
	Port                            string
	PortHttp                        string

	MariaDBHost                     string
	MariaDBPort                     string
	MariaDBUser                     string
	MariaDBPassword                 string
	MariaDBDatabase                 string

	AuthGitHubClientID              string
	AuthGitHubClientSecret          string
	AuthGitHubClientCallbackURL     string
	
	FrontendURL                     string
	FrontendURLAuthCallback         string

	CookieDomain                    string
	CookiePath                      string
	CookieHTTPOnly                  bool
	CookieSecure                    bool
	// CookieSameSite               string
	// CookieMaxAge                 int
	// CookieExpire                 int

	SSLKey                          string
	SSLCert                         string

	JWTSecret                       string
	ServiceToken                    string
}

func LoadConfig() *Config {
	// cookieMaxAge, err := strconv.Atoi(utils.GetEnv("COOKIE_MAX_AGE", "86400"))
	// if err != nil {
	// 	cookieMaxAge = 86400
	// }

	// cookieExpire, _ := strconv.Atoi(utils.GetEnv("COOKIE_EXPIRE", "86400"))
	// if err != nil {
	// 	cookieExpire = 86400
	// }

	return &Config{
		Host:                            utils.GetEnv("HOST", "localhost"),
		Port:                            utils.GetEnv("PORT", "5000"),
		PortHttp:                        utils.GetEnv("PORT_HTTP", "5001"),

		MariaDBHost:                     utils.GetEnv("MARIADB_HOST", "localhost"),
		MariaDBPort:                     utils.GetEnv("MARIADB_PORT", "3306"),
		MariaDBUser:                     utils.GetEnv("MARIADB_USER", "root"),
		MariaDBPassword:                 utils.GetEnv("MARIADB_PASSWORD", "root"),
		MariaDBDatabase:                 utils.GetEnv("MARIADB_DATABASE", "codeduel"),

		AuthGitHubClientID:              utils.GetEnv("AUTH_GITHUB_CLIENT_ID", ""),
		AuthGitHubClientSecret:          utils.GetEnv("AUTH_GITHUB_CLIENT_SECRET", ""),
		AuthGitHubClientCallbackURL:     utils.GetEnv("AUTH_GITHUB_CLIENT_CALLBACK_URL", "http://localhost:5000/auth/github/callback"),
		
		FrontendURL:                     utils.GetEnv("FRONTEND_URL", "http://localhost:5173"),
		FrontendURLAuthCallback:         utils.GetEnv("FRONTEND_URL_AUTH_CALLBACK", "http://localhost:5173/auth/callback"),

		CookieDomain:                    utils.GetEnv("COOKIE_DOMAIN", "localhost"),
		CookiePath:                      utils.GetEnv("COOKIE_PATH", "/"),
		CookieHTTPOnly:                  utils.GetEnv("COOKIE_HTTP_ONLY", "true") == "true",
		CookieSecure:                    utils.GetEnv("COOKIE_SECURE", "false") == "true",
		// CookieSameSite:                  utils.GetEnv("COOKIE_SAME_SITE", "Lax"),
		// CookieMaxAge:                    cookieMaxAge,
		// CookieExpire:                    cookieExpire,

		SSLKey:                          utils.GetEnv("SSL_KEY", "ssl/server.key"),
		SSLCert:                         utils.GetEnv("SSL_CERT", "ssl/server.crt"),

		JWTSecret:                       utils.GetEnv("JWT_SECRET", "yeahSuperSecret"),
		ServiceToken:                    utils.GetEnv("SERVICE_TOKEN", "yeahSuperToken"),
	}
}
