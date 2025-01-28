//go:build example
// +build example

package examples

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/goccy/go-json"

	"shortly/internal/app/config"
	"shortly/internal/app/dto"
	"shortly/internal/app/repository"
	"shortly/internal/app/router"
	"shortly/internal/app/worker"
	"shortly/internal/logger"
)

func Example_handleCreateShortLink() {
	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	appLogger := logger.NewLogger()
	appRepo := repository.NewInMemoryRepository()
	deleteWorker := worker.NewDeleteWorker(ctx, cfg, appRepo, appLogger)
	deleteWorker.Start()
	defer deleteWorker.Stop()

	appRouter := router.NewRouter(cfg, appRepo, deleteWorker, appLogger)
	appServer := httptest.NewServer(appRouter)
	defer appServer.Close()

	body := `{"url":"https://github.com"}`
	resp, err := http.Post(appServer.URL+"/api/shorten", "application/json", strings.NewReader(body))
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Body:", string(b))
	// Output:
	// Status: 201
	// Body: {"result":"http://localhost:8080/abcd1234"}
}

func Example_handleGetShortLink() {
	ctx := context.Background()
	cfg := &config.Config{
		BaseURL: "http://localhost:8080",
	}
	appLogger := logger.NewLogger()
	appRepo := repository.NewInMemoryRepository()
	deleteWorker := worker.NewDeleteWorker(ctx, cfg, appRepo, appLogger)
	deleteWorker.Start()
	defer deleteWorker.Stop()

	appRouter := router.NewRouter(cfg, appRepo, deleteWorker, appLogger)
	appServer := httptest.NewServer(appRouter)
	defer appServer.Close()

	createBody := `{"url":"https://golang.org"}`
	respCreate, _ := http.Post(appServer.URL+"/api/shorten", "application/json", strings.NewReader(createBody))
	defer respCreate.Body.Close()
	cBytes, _ := io.ReadAll(respCreate.Body)

	var created dto.CreateShortLinkResponse
	_ = json.Unmarshal(cBytes, &created)
	shortCode := strings.TrimPrefix(created.Result, cfg.BaseURL+"/")

	resp, err := http.Get(appServer.URL + "/api/shorten/" + shortCode)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println("Status:", resp.StatusCode)
	fmt.Println("Body:", string(b))
	// Output:
	// Status: 200
	// Body: {"result":"https://golang.org"}
}
