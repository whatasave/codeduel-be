package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     string
	PortHttp string

	MariaDBHost     string
	MariaDBPort     string
	MariaDBUser     string
	MariaDBPassword string
	MariaDBDatabase string

	AuthGitHubClientID          string
	AuthGitHubClientSecret      string
	AuthGitHubClientCallbackURL string

	FrontendURL string

	CookieDomain   string
	CookiePath     string
	CookieHTTPOnly bool
	CookieSecure   bool
	// CookieSameSite               string
	// CookieMaxAge                 int
	// CookieExpire                 int

	CorsOrigin      string
	CorsMethods     string
	CorsHeaders     string
	CorsCredentials bool

	SSLKey  string
	SSLCert string

	JWTSecret                       string
	JWTExpiresInMinutes             int
	JWTRefreshTokenExpiresInMinutes int

	ServiceToken string
}

var config *Config

func LoadConfig() *Config {
	if config == nil {
		// loading env only if not in production
		if GetEnv("ENV", "development") == "development" {
			if err := godotenv.Load(); err != nil {
				log.Printf("%s Error loading .env file", GetLogTag("ENV"))
			}
		}

		config = &Config{
			Host:     GetEnv("HOST", "localhost"),
			Port:     GetEnv("PORT", "5000"),
			PortHttp: GetEnv("PORT_HTTP", "5001"),

			MariaDBHost:     GetEnv("MARIADB_HOST", "localhost"),
			MariaDBPort:     GetEnv("MARIADB_PORT", "3306"),
			MariaDBUser:     GetEnv("MARIADB_USER", "codeduel"),
			MariaDBPassword: GetEnv("MARIADB_PASSWORD", "codeduel"),
			MariaDBDatabase: GetEnv("MARIADB_DATABASE", "codeduel"),

			AuthGitHubClientID:          GetEnv("AUTH_GITHUB_CLIENT_ID", ""),
			AuthGitHubClientSecret:      GetEnv("AUTH_GITHUB_CLIENT_SECRET", ""),
			AuthGitHubClientCallbackURL: GetEnv("AUTH_GITHUB_CLIENT_CALLBACK_URL", "http://localhost:5000/auth/github/callback"),

			FrontendURL: GetEnv("FRONTEND_URL", "http://localhost:5173"),

			CookieDomain:   GetEnv("COOKIE_DOMAIN", "localhost"),
			CookiePath:     GetEnv("COOKIE_PATH", "/"),
			CookieHTTPOnly: GetEnv("COOKIE_HTTP_ONLY", "true") == "true",
			CookieSecure:   GetEnv("COOKIE_SECURE", "false") == "true",
			// CookieSameSite:                  GetEnv("COOKIE_SAME_SITE", "Lax"),
			// CookieMaxAge:                    toInt(GetEnv("COOKIE_MAX_AGE", "86400"), 60 * 60 * 24), // 24 hours
			// CookieExpire:                    toInt(GetEnv("COOKIE_EXPIRE", "86400"), 60 * 60 * 24), // 24 hours

			CorsOrigin:      GetEnv("CORS_ORIGIN", "http://localhost:5173"),
			CorsMethods:     GetEnv("CORS_METHODS", "GET, POST, PUT, PATCH, DELETE, OPTIONS"),
			CorsHeaders:     GetEnv("CORS_HEADERS", "Content-Type, x-token, Accept, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token"),
			CorsCredentials: GetEnv("CORS_CREDENTIALS", "true") == "true",

			SSLKey:  GetEnv("SSL_KEY", "ssl/server.key"),
			SSLCert: GetEnv("SSL_CERT", "ssl/server.crt"),

			JWTSecret:                       GetEnv("JWT_SECRET", GenerateRandomString(16)),
			JWTExpiresInMinutes:             ToInt(GetEnv("JWT_EXPIRES_IN_MINUTES", "5"), 5),
			JWTRefreshTokenExpiresInMinutes: ToInt(GetEnv("JWT_REFRESH_TOKEN_EXPIRES_IN_MINUTES", "43200"), 60*24*30), // 30 days

			ServiceToken: GetEnv("SERVICE_TOKEN", "yeahSuperToken"),
		}
	}

	return config
}

func GetEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Printf("%s%s %s variable not found, using default value: %s\n", GetLogTag("ENV"), GetLogTag("warn"), key, defaultValue)
		return defaultValue
	}

	return value
}
