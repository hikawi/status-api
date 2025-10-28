package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"luny.dev/status-api/services"
)

type GetServicesResponse struct {
	Result    map[string]services.CheckResult `json:"result"`
	CheckedAt int64                           `json:"checked_at_milli"`
}

var (
	lastResult GetServicesResponse
	cacheMutex sync.RWMutex
)

const cacheDurationMilli = 5 * 60 * 1000

func GetServicesHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	now := time.Now()

	cacheMutex.RLock()
	isFresh := now.UnixMilli()-lastResult.CheckedAt <= cacheDurationMilli
	cachedResult := lastResult
	cacheMutex.RUnlock()

	// Cache hit
	if isFresh {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(cachedResult); err != nil {
			log.Printf("error: encoding json response: %v", err)
		}
		return
	}

	// Cache miss
	// Acquire a Write Lock before modifying the shared variable
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// Double-check: Another request might have updated the cache while we were waiting for the lock
	if now.UnixMilli()-lastResult.CheckedAt < cacheDurationMilli {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(lastResult); err != nil {
			log.Printf("error: encoding json response: %v", err)
		}
		return
	}

	// Perform the slow query
	newResult := services.QueryAllServices()
	newResponse := GetServicesResponse{newResult, now.UnixMilli()}

	// Update the cache
	lastResult = newResponse

	// --- 3. Return New Data ---
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(newResponse); err != nil {
		log.Printf("error: encoding json response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
