package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/v4/stdlib" // Import the pgx driver
)

const webPort = "80"

func main() {
	// Connect to database
	db := initDB()

	// Create sessions
	session := initSession()

	// Create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Create waitgroup
	wg := sync.WaitGroup{}

	// Set up application server
	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &wg,
	}

	// Listen for web connection
	app.serve()
}

func (app *Config) serve() {
	// Start HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("Starting web server...")

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func initDB() *sql.DB {
	// Call another function to connect to database with retries
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to database")
	}
	return conn
}

func connectToDB() *sql.DB {
	counts := 0
	maxRetries := 10
	// Connection string
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet:", err)
		} else {
			log.Println("Connected to database")
			return connection
		}

		if counts >= maxRetries {
			log.Println("Exceeded maximum number of retries")
			return nil
		}

		counts++
		log.Printf("Backing off for 1 second (%d/%d)\n", counts, maxRetries)
		time.Sleep(1 * time.Second)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Initialize the session manager with Redis store
func initSession() *scs.SessionManager {
	session := scs.New()
	session.Store = redisstore.New(initRedis()) // Use Redis store
	session.Lifetime = 24 * time.Hour           // Set session lifetime to 24 hours
	session.Cookie.Persist = true               // Make session cookie persistent
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = true // Cookie is only sent over HTTPS

	return session
}

// Initialize the Redis connection pool
func initRedis() *redis.Pool {
	redisPool := &redis.Pool{
		MaxIdle: 10, // Maximum number of idle connections
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", os.Getenv("REDIS")) // Get Redis address from environment variable
		},
	}
	return redisPool
}
