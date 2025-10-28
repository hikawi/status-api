package routes

import (
	"encoding/json"
	"log"
	"net/http"
)

type GetRootResponse struct {
	Message string `json:"message"`
}

func GetRootHandler(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	data := GetRootResponse{
		Message: "Hello, World!",
	}

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("error: encoding json response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
