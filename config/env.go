package config

import "github.com/xedom/codeduel/utils"

type Config struct {
	Host                            string
	Port                            string   

	MariaDBHost                     string
	MariaDBPort                     string   
	MariaDBUser                     string
	MariaDBPassword                 string
	MariaDBDatabase                 string
	
	FrontendURL                     string
	FrontendURLAuthCallback         string

	AuthGitHubClientID              string
	AuthGitHubClientSecret          string
	AuthGitHubClientCallbackURL     string
}

func LoadConfig() *Config {
	return &Config{
		Host:                            utils.GetEnv("HOST", "localhost"),
		Port:                            utils.GetEnv("PORT", "5000"),

		MariaDBHost:                     utils.GetEnv("MARIADB_HOST", "localhost"),
		MariaDBPort:                     utils.GetEnv("MARIADB_PORT", "3306"),
		MariaDBUser:                     utils.GetEnv("MARIADB_USER", "root"),
		MariaDBPassword:                 utils.GetEnv("MARIADB_PASSWORD", "root"),
		MariaDBDatabase:                 utils.GetEnv("MARIADB_DATABASE", "codeduel"),
		
		FrontendURL:                     utils.GetEnv("FRONTEND_URL", "http://localhost:5173"),
		FrontendURLAuthCallback:         utils.GetEnv("FRONTEND_URL_AUTH_CALLBACK", "http://localhost:5173/auth/callback"),

		AuthGitHubClientID:              utils.GetEnv("AUTH_GITHUB_CLIENT_ID", ""),
		AuthGitHubClientSecret:          utils.GetEnv("AUTH_GITHUB_CLIENT_SECRET", ""),
		AuthGitHubClientCallbackURL:     utils.GetEnv("AUTH_GITHUB_CLIENT_CALLBACK_URL", "http://localhost:8080/auth/github/callback"),
	}
}
