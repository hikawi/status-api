// Package services provides a set of functions that are independent from endpoints
package services

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Service struct {
	ID  string
	URL string
}

type CheckResult struct {
	ID           string `json:"id"`
	Status       string `json:"status"`
	ResponseCode int    `json:"response_code"`
	ResponseTime int64  `json:"response_time"`
}

func GetAllServices() []Service {
	services := make([]Service, 0, 13)

	services = append(services, Service{"online-news", "https://online-news.luny.dev"})
	services = append(services, Service{"it-tools", "https://it-tools.luny.dev"})
	services = append(services, Service{"lubook", "https://lubook.luny.dev"})
	services = append(services, Service{"lubook-api", "https://api.lubook.luny.dev"})
	services = append(services, Service{"acci-api", "https://api.acci.luny.dev/health"})

	services = append(services, Service{"acci", "https://acci.luny.dev"})
	services = append(services, Service{"cc-sakura", "https://cc-sakura.luny.dev"})
	services = append(services, Service{"frilly-dev", "https://www.frilly.dev"})
	services = append(services, Service{"locket-api", "https://locket.luny.dev"})
	services = append(services, Service{"pomodoro", "https://pomodoro.luny.dev"})

	services = append(services, Service{"tic-tac-toe", "https://tic-tac-toe.luny.dev"})
	services = append(services, Service{"status", "https://status.luny.dev/health"})
	services = append(services, Service{"sakila-api", "https://sakila.luny.dev/api/health"})

	return services
}

func CheckURL(id string, url string, results chan<- CheckResult, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("info: service %s: checking %s", id, url)

	client := http.Client{Timeout: 5 * time.Second}

	start := time.Now().UnixMilli()
	res, err := client.Get(url)
	responseTimeMs := time.Now().UnixMilli() - start

	var checkResult CheckResult
	checkResult.ID = id

	if err != nil {
		fmt.Printf("error: can't connect to %s for service %s\n", url, id)
		checkResult.Status = "down"
		checkResult.ResponseTime = -1
	} else {
		fmt.Printf("success: connected to %s for service %s\n", url, id)
		checkResult.ResponseCode = res.StatusCode
		checkResult.ResponseTime = responseTimeMs

		if res.StatusCode >= 500 {
			checkResult.Status = "down"
		} else if res.StatusCode/100 == 2 {
			checkResult.Status = "up"
		} else {
			checkResult.Status = "unknown"
		}

		io.Copy(io.Discard, res.Body)
	}

	results <- checkResult
}

func QueryAllServices() map[string]CheckResult {
	services := make(map[string]CheckResult)
	allServices := GetAllServices()
	fmt.Printf("info: querying for all services (%d)\n", len(allServices))

	var wg sync.WaitGroup
	results := make(chan CheckResult, len(allServices))

	for _, service := range allServices {
		wg.Add(1)
		go CheckURL(service.ID, service.URL, results, &wg)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		services[result.ID] = result
	}

	return services
}
