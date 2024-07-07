package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/arthben/http_jwt_crud/api/handlers"
	"github.com/arthben/http_jwt_crud/api/middlewares"
	"github.com/arthben/http_jwt_crud/internal/config"
	"github.com/arthben/http_jwt_crud/internal/database"
	"github.com/arthben/http_jwt_crud/internal/logging"
)

type App struct {
	cfg    *config.EnvParams
	logger *slog.Logger
}

func main() {

	// load config
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("err load config: %v\n", err)
		os.Exit(1)
		return
	}

	// init slog
	logger := logging.NewLogger(cfg.AppMode, "my-service")

	// init database
	db, err := database.BoostrapDatabase(&cfg)
	if err != nil {
		logger.Error("err Open Database:", slog.String("", err.Error()))
		os.Exit(1)
		return
	}
	defer database.CloseDatabase(db)

	handler, err := handlers.NewHandlers(db, &cfg).BuildRouter()
	if err != nil {
		logger.Error("err BuildRouter:", slog.String("", err.Error()))
		os.Exit(1)
		return
	}

	// init middleware
	middlewareChain := middlewares.MiddlewareChain(
		middlewares.CORS,
		middlewares.ContextRequestID,
		middlewares.RequestLogger,
	)

	timeout, _ := strconv.Atoi(cfg.ServerTimeout)
	server := http.Server{
		Addr:              "0.0.0.0:" + cfg.Port,
		Handler:           middlewareChain(logger, handler),
		ReadTimeout:       time.Duration(timeout) * time.Second,
		ReadHeaderTimeout: time.Duration(timeout) * time.Second,
		WriteTimeout:      time.Duration(timeout) * time.Second,
	}
	defer server.Close()

	osSignal := make(chan os.Signal, 2)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("SERVER START", slog.String("address", "0.0.0.0:"+cfg.Port))
		logger.Warn("SERVER CLOSE", slog.String("error", server.ListenAndServe().Error()))
	}()

	// wait until server closed
	select {
	case <-osSignal:
		logger.Warn("Server Closed Interrupted by OS")
	}

	logger.Info("SERVER STOP")
}
