package main

import (
	"net/http"

	"github.com/rs/cors"
	"luny.dev/status-api/routes"
)

func main() {
	mux := http.NewServeMux()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
		Debug:            true,
	})

	mux.HandleFunc("/health", routes.GetHealthHandler)
	mux.HandleFunc("/services", routes.GetServicesHandler)

	handler := c.Handler(mux)

	http.ListenAndServe(":80", handler)
}
