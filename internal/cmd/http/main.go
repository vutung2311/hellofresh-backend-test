package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"vutung2311-golang-test/internal/handler"
	"vutung2311-golang-test/internal/repository"
	"vutung2311-golang-test/pkg/httpclient"
	"vutung2311-golang-test/pkg/router"
	"vutung2311-golang-test/pkg/router/middleware"
	"vutung2311-golang-test/pkg/tracing"
	"vutung2311-golang-test/pkg/worker"
)

func main() {
	baseURL := os.Getenv("BASE_URL")
	if len(baseURL) == 0 {
		panic("BASE_URL is required")
	}
	logger := logrus.New()
	logger.SetFormatter(new(logrus.JSONFormatter))
	loggerCreator := func(ctx context.Context) logrus.FieldLogger {
		return logger.WithField("contextId", tracing.GetReqID(ctx))
	}
	httpClient := httpclient.New(10 * time.Second).WithRequestResponseLogger(loggerCreator)
	workerPool := worker.NewPool(1000, loggerCreator)
	recipeGetHandler := handler.CreateRecipeGetHandler(
		repository.NewRecipeRepository(baseURL, httpClient, workerPool),
	)
	r := router.New(middleware.RequestResponseLogger(logger), middleware.RequestID)
	r.AddRoute("/recipes", recipeGetHandler)

	server := http.Server{
		Handler:           r,
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       time.Minute,
	}

	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
