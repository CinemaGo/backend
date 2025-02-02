package main

import "cinemaGo/backend/api/routes"

func main() {
	router := routes.Router()

	router.Run()
}
