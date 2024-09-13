package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

const webport = "90"

func main() {
	//connect to database
	db := initDB()

	//create sessions
	session := initSession()

	//create loggers.....
	info_Log := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	error_Log := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime)

	//create channel

	//create waitgroup
	wg := sync.WaitGroup{}

	//set up application server
	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  info_Log,
		ErrorLog: error_Log,
		Wait:     &wg,
	}

	//set up mail

	//listen for web connection

}
func initDB() *sql.DB {
	// i will call another function just because
	// i can try to connect to database repeatedly if necessary

	conn := connectToDB()
	if conn == nil {
		log.Panic("cant connect to database")

	}

	return nil
}
func connectToDB() *sql.DB {
	counts := 0
	//i will try to connect to database fixed number of times
	//if it fails after that time it dies

	//connection string
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Panicln("Posgtres not ready yet")

		} else {
			log.Println("connected to database")
			return connection
		}
		if counts > 10 {
			return nil
		}
		log.Print("Backing off for 1 second")
		time.Sleep(time.Second * 1)
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
	//it Configures the session cookie to persist in the user's browser even after the browser is closed.
	session.Cookie.SameSite = http.SameSiteLaxMode
	//Sets the SameSite attribute to Lax, which allows the cookie to be sent with same-site requests and some cross-site requests.
	session.Cookie.Secure = true // Cookie is only sent over HTTPS
	// Ensures the cookie is only sent over HTTPS

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
