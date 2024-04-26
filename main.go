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
			log.Printf("%s Error loading .env file", utils.GetLogTag("main"))
		}
	}
	loadConfig := config.LoadConfig()

	mariaDB, err := db.NewDB(
		loadConfig.MariaDBHost,
		loadConfig.MariaDBPort,
		loadConfig.MariaDBUser,
		loadConfig.MariaDBPassword,
		loadConfig.MariaDBDatabase,
	)
	if err != nil {
		log.Printf("%s%s Error creating DB instance: %v", utils.GetLogTag("DB"), utils.GetLogTag("error"), err.Error())
		return
	}
	if mariaDB == nil {
		log.Printf("%s%s Error creating DB instance: %v", utils.GetLogTag("DB"), utils.GetLogTag("error"), "DB instance is nil")
		return
	}
	if err := mariaDB.InitUserTables(); err != nil {
		log.Printf("%s%s Error initializing DB user tables: %v", utils.GetLogTag("DB"), utils.GetLogTag("error"), err.Error())
	}
	if err := mariaDB.InitMatchTables(); err != nil {
		log.Printf("%s%s Error initializing DB match tables: %v", utils.GetLogTag("DB"), utils.GetLogTag("error"), err.Error())
	}
	if err := mariaDB.InitLobbyTables(); err != nil {
		log.Printf("%s%s Error initializing DB lobby tables: %v", utils.GetLogTag("DB"), utils.GetLogTag("error"), err.Error())
	}

	server := api.NewAPIServer(loadConfig, mariaDB)
	err = server.Run()
	if err != nil {
		log.Printf("%s Error running server: %v", utils.GetLogTag("main"), err.Error())
	}

	err = mariaDB.Close()
	if err != nil {
		log.Printf("%s Error closing DB connection: %v", utils.GetLogTag("DB"), err.Error())
	}
}
