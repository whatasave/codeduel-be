package main

import (
	"log"

	"github.com/xedom/codeduel/api"
	"github.com/xedom/codeduel/db"
	"github.com/xedom/codeduel/utils"
)

func main() {
	loadConfig := utils.LoadConfig()

	mariaDB, err := db.NewDB(
		loadConfig.MariaDBHost,
		loadConfig.MariaDBPort,
		loadConfig.MariaDBUser,
		loadConfig.MariaDBPassword,
		loadConfig.MariaDBDatabase,
	)

	if err != nil || mariaDB == nil {
		log.Printf("%s%s Error creating DB instance: %v", utils.GetLogTag("DB"), utils.GetLogTag("error"), "DB instance is nil")
		return
	}

	if err := mariaDB.MigrationBulk(
		mariaDB.InitUserTables(),
		mariaDB.InitLobbyTables(),
		mariaDB.InitChallengeTables(),
	); err != nil {
		log.Printf("%s%s Error migrating DB user tables: %v", utils.GetLogTag("DB"), utils.GetLogTag("error"), err.Error())
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
