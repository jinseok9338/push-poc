package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

func main() {
	// Get the database URL from the environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable not set")
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a new echo instance
	e := echo.New()

	// Define your API routes and handlers here
	// make hello world handler
	helloWorld := func(c echo.Context) error {
		return c.String(200, "Hello, !!!")
	}

	// Start the server with handler
	e.GET("/", helloWorld)
	e.Logger.Fatal(e.Start(":3000"))

}
