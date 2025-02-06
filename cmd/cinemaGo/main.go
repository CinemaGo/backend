package main

import (
	"cinemaGo/backend/api/handlers"
	"cinemaGo/backend/api/routes"
	"cinemaGo/backend/internal/models"
	"cinemaGo/backend/internal/services"
	"cinemaGo/backend/pkg/configs"
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

// main is the entry point of the application.
// It sets up the environment, connects to the database, and starts the HTTP server.
func main() {
	// Load environment variables from the .env file.
	// If the file is not found or there's an error, the program will terminate.
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Load the database username from the environment variables.
	// If the variable is missing or there's an error, the program will terminate.
	dbUserName, err := configs.LoadEnvironmentVariable("DB_USER_NAME")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Load the database password from the environment variables.
	// If the variable is missing or there's an error, the program will terminate.
	dbUserPassword, err := configs.LoadEnvironmentVariable("DB_USER_PASSWORD")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Load the database name from the environment variables.
	// If the variable is missing or there's an error, the program will terminate.
	dbName, err := configs.LoadEnvironmentVariable("DB_NAME")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Load the database SSL mode setting from the environment variables.
	// If the variable is missing or there's an error, the program will terminate.
	dbSSLMode, err := configs.LoadEnvironmentVariable("DB_SSL_MODE")
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Construct the database connection string using the loaded environment variables.
	// The connection string includes the user, password, database name, and SSL mode.
	dsn := fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v", dbUserName, dbUserPassword, dbName, dbSSLMode)

	// Open a connection to the database using the connection string.
	// If the connection fails, the program will terminate with an error message.
	db, err := models.OpenDB(dsn)
	if err != nil {
		log.Fatalf("%v", err)
	}

	// Ensure the database connection is closed when the program exits.
	defer db.Close()

	moviesService, err := services.NewMoviesService(db)
	if err != nil {
		log.Fatal(err)
	}
	moviesHandler := handlers.NewMoviesHandler(moviesService)

	usersService, err := services.NewUsersService(db)
	if err != nil {
		log.Fatal(err)
	}
	usersHandler := handlers.NewUsersHandler(usersService)

	bookingService := services.NewBookingService(db)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	serveHandlersWrapper := routes.ServeHandlersWrapper{
		MoviesHandler:  moviesHandler,
		UsersHandler:   usersHandler,
		BookingHandler: bookingHandler,
	}

	router := routes.Router(&serveHandlersWrapper)

	router.Run()
}
