package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/xedom/codeduel/types"
	"github.com/xedom/codeduel/utils"
)

type DB interface {
	GetUsers() ([]*types.UserResponse, error)
	GetUserByID(int) (*types.User, error)
	GetUserByUsername(string) (*types.User, error)
	GetUserStats(int) ([]*types.UserStatsParsed, error)
	CreateUser(*types.User) error
	UpdateUser(*types.User) error
	DeleteUser(int) error
	DeleteUserByUsername(string) error

	GetChallenges() (*[]types.Challenge, error)
	GetChallengeByID(int) (*types.Challenge, error)
	CreateChallenge(*types.Challenge) error
	UpdateChallenge(*types.Challenge) error
	DeleteChallenge(int) error
	GetChallengesByOwnerID(int) (*[]types.Challenge, error)

	CreateLobby(*types.Lobby) error
	CreateLobbyUser(int, int) error
	UpdateLobbyUserSubmission(*types.LobbyUser) error
	GetLobbyByUniqueId(string) (*types.Lobby, error)
	GetLobbyResults(string) (*types.LobbyResults, error)
	EndLobby(string) error
	UpdateShareLobbyCode(int, int, bool) error
	GetMatchByUsername(string) ([]*types.SingleMatchResult, error)

	GetAuthByProviderAndID(string, string) (*types.AuthEntry, error)
	CreateAuth(*types.AuthEntry) error
	CreateRefreshToken(int, *utils.JWT) error
	DeleteRefreshToken(string) error
}

type MariaDB struct {
	db *sql.DB
}

type MigrationFunc func() error

func NewDB(host, port, user, pass, name string) (*MariaDB, error) {
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + name

	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// pool.SetConnMaxLifetime(0)
	// pool.SetMaxIdleConns(15)
	// pool.SetMaxOpenConns(15)

	if err := pool.Ping(); err != nil {
		return nil, err
	}

	var version string
	err = pool.QueryRow("SELECT VERSION();").Scan(&version)
	if err != nil {
		return nil, err
	}
	log.Printf("%s Connected to: %s", utils.GetLogTag("db"), version)

	return &MariaDB{db: pool}, nil
}

var pool *sql.DB

func InitOLD() {
	dbUser := os.Getenv("MARIADB_USER")
	dbPass := os.Getenv("MARIADB_PASSWORD")
	dbName := os.Getenv("MARIADB_DATABASE")
	var err error

	dbConnectionUri := dbUser + ":" + dbPass + "@/" + dbName
	pool, err = sql.Open("mysql", dbConnectionUri)
	if err != nil {
		log.Fatal("[DB] Unable to connect to the database:", err)
	}

	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	appSignal := make(chan os.Signal, 3)
	signal.Notify(appSignal, os.Interrupt)

	go func() {
		<-appSignal
		stop()
	}()

	Ping(ctx)

	var version string
	pool.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("[DB] Connected to:", version)
	// Query(ctx, *id)

	err = pool.Close()
	if err != nil {
		log.Fatal("[DB] Unable to close the database connection:", err)
	}
}

func Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("[DB] Unable to connect to the database: %v", err)
	}
}

func (m *MariaDB) InitDatabase() error {
	query := "CREATE DATABASE IF NOT EXISTS `codeduel`;"
	_, err := m.db.Exec(query)

	return err
}

func (m *MariaDB) Migration(migrationFuncs ...MigrationFunc) error {
	for _, migration := range migrationFuncs {
		if err := migration(); err != nil {
			return err
		}
	}

	return nil
}

func (m *MariaDB) MigrationBulk(migrationFuncs ...[]MigrationFunc) error {
	for _, migrations := range migrationFuncs {
		// migrations = append(migrations, migrations...)
		err := m.Migration(migrations...)
		if err != nil {
			return err
		}
	}

	return nil
}

// func Query(ctx context.Context, id int64) {
// 	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
// 	defer cancel()

// 	var name string
// 	err := pool.QueryRowContext(ctx, "select p.name from people as p where p.id = :id;", sql.Named("id", id)).Scan(&name)
// 	if err != nil {
// 		log.Fatal("[DB] unable to execute search query", err)
// 	}
// 	log.Println("name=", name)
// }

func (m *MariaDB) Close() error {
	return m.db.Close()
}
