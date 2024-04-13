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
  if utils.GetEnv("ENV", "development") == "development" {
    if err := godotenv.Load(); err != nil {
      log.Println("%s Error loading .env file", utils.GetLogTag("main"))
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
  if err != nil { log.Printf("%s Error creating DB instance: %v", utils.GetLogTag("main"), err) }
  defer db.Close()
  if err := db.InitUserTables(); err != nil { log.Printf("%s Error initializing DB user tables: %v", utils.GetLogTag("main"), err) }
  if err := db.InitMatchTables(); err != nil { log.Printf("%s Error initializing DB match tables: %v", utils.GetLogTag("main"), err) }

  server := api.NewAPIServer(config, db)
  server.Run()
}
