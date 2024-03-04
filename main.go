package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/xedom/codeduel/api"
	"github.com/xedom/codeduel/config"
	"github.com/xedom/codeduel/db"
	"github.com/xedom/codeduel/utils"
)


func main() {
  // loading env only if not in production
  if utils.GetEnv("ENV", "development") != "production" {
    if err := godotenv.Load(); err != nil {
      log.Println("[MAIN] Error loading .env file")
    }
  }
  config := config.LoadConfig()

  db, err := db.NewDB(
    config.MariaDBHost,
    config.MariaDBPort,
    config.MariaDBUser,
    config.MariaDBPassword,
    config.MariaDBDatabase,
  )
  if err != nil { log.Printf("[MAIN] Error creating DB instance: %v", err) }
  defer db.Close()
  if err := db.InitUserTables(); err != nil { log.Printf("[MAIN] Error initializing DB user tables: %v", err) }
  if err := db.InitMatchTables(); err != nil { log.Printf("[MAIN] Error initializing DB match tables: %v", err) }

  server := api.NewAPIServer(config, db)
  server.Run()
}
