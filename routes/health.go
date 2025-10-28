// Package routes hosts a list of routes for the Status API.
package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

type GetHealthResponse struct {
	Status string `json:"status"`
}

func GetHealthHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(rw, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	rw.Header().Add("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)

	data := GetHealthResponse{
		Status: "healthy",
	}

	if err := json.NewEncoder(rw).Encode(data); err != nil {
		log.Printf("error: encoding json response: %v", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
