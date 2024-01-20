package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/xedom/codeduel/api"
	"github.com/xedom/codeduel/db"
)

func main() {
  loadingEnvVars()
  warnUndefinedEnvVars()
  
  db, err := db.NewDB(
    os.Getenv("MARIADB_HOST"),
    os.Getenv("MARIADB_PORT"),
    os.Getenv("MARIADB_USER"),
    os.Getenv("MARIADB_PASSWORD"),
    os.Getenv("MARIADB_DATABASE"),
  )
  if err != nil { log.Printf("[MAIN] Error creating DB instance: %v", err) }
  defer db.Close()
  if err := db.InitUserTables(); err != nil { log.Printf("[MAIN] Error initializing DB user tables: %v", err) }
  if err := db.InitMatchTables(); err != nil { log.Printf("[MAIN] Error initializing DB match tables: %v", err) }

  server := api.NewAPIServer(os.Getenv("HOST"), os.Getenv("PORT"), db)
  server.Run()
}

func loadingEnvVars() {
  isProduction := os.Getenv("GO_ENV") == "production"
  if isProduction { return }
  
  path_dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
  if err != nil { log.Printf("[MAIN] Error getting absolute path: %v", err) }
  // log.Default().SetPrefix("[MAIN] ")
  log.Printf("Loading .env file from %s", path_dir)
  env_path := filepath.Join(path_dir, ".env")
  if _, err := os.Stat(env_path); os.IsNotExist(err) {
    log.Printf("[MAIN] Error: .env file not found in %s", path_dir)
    return
  }
  err = godotenv.Load(env_path)
  if err != nil { log.Printf("[MAIN] Error loading .env file") }
}

func warnUndefinedEnvVars() {
  envVars := []string{
    "HOST",
    "PORT",
    "FRONTEND_URL",
    "FRONTEND_URL_AUTH_CALLBACK",
    "MARIADB_HOST",
    "MARIADB_PORT",
    "MARIADB_USER",
    "MARIADB_PASSWORD",
    "MARIADB_DATABASE",
    "AUTH_GITHUB_CLIENT_ID",
    "AUTH_GITHUB_CLIENT_SECRET",
    "AUTH_GITHUB_CLIENT_CALLBACK_URL",
  }
  
  for _, envVar := range envVars {
    test, exists := os.LookupEnv(envVar)
    if !exists { log.Printf("[MAIN] Warning: %s not defined in .env file", envVar) }
    if test == "" { log.Printf("[MAIN] Warning: %s is empty", envVar) }
    log.Printf("[MAIN] %s: %s", envVar, test)
  }
}