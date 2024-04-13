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
	CreateUser(*types.User) error
	DeleteUser(int) error
	UpdateUser(*types.User) error
	GetUsers() ([]*types.User, error)
	GetUserByID(int) (*types.User, error)

	CreateChallenge(*types.Challenge) error
	DeleteChallenge(int) error
	UpdateChallenge(*types.Challenge) error
	GetChallenges() (*[]types.Challenge, error)
	GetChallengeByID(int) (*types.Challenge, error)

	CreateAuth(*types.AuthEntry) error
	GetAuthByProviderAndID(string, string) (*types.AuthEntry, error)
}

type MariaDB struct {
	db *sql.DB
}

func NewDB(host, port, user, pass, name string) (*MariaDB, error) {
	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + name
	pool, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	pool.SetConnMaxLifetime(0)
	pool.SetMaxIdleConns(5)
	pool.SetMaxOpenConns(5)

	if err := pool.Ping(); err != nil {
		return nil, err
	}

	var version string
	pool.QueryRow("SELECT VERSION()").Scan(&version)
	log.Printf("%s Connected to: %s", utils.GetLogTag("db"), version)

	return &MariaDB{
		db: pool,
	}, nil
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
	defer pool.Close()

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
}

func Ping(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err := pool.PingContext(ctx); err != nil {
		log.Fatalf("[DB] Unable to connect to the database: %v", err)
	}
}

func (m *MariaDB) InitDatabase() error {
	query := `CREATE DATABASE IF NOT EXISTS codeduel;`
	_, err := m.db.Exec(query)

	return err
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
