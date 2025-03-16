package driver

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq" // Import driver PostgreSQL
)

var db *sql.DB

var rd *redis.Client

// StartUpDB initializes the database connection
func StartUpDB() {
	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	fmt.Println("Waiting for the database to start up...")

	time.Sleep(5 * time.Second) // Delay for the database to start up

	// Open database connection
	var err error
	db, err = sql.Open("postgres", connStr) // Correct driver name and variable assignment

	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Ping database to verify connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	fmt.Println("Successfully connected to the database")

}

func InitRedis() {
	ctx := context.Background()

	rd = redis.NewClient(&redis.Options{
		Addr:     "192.168.249.100:6379",
		Password: "Ocb1234*",
		DB:       0,
	})

	pong, err := rd.Ping(ctx).Result()
	fmt.Println("Successfully connected to redis database: ", pong, err)

}

func GetRedis() *redis.Client {
	return rd
}

// GetDBCarManageMent returns the database connection
func GetDBCarManageMent() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() {
	if err := db.Close(); err != nil {
		log.Fatalf("Error closing the database: %v", err)
	}
}
