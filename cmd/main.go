package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"deu/internal/config"
	"deu/internal/places"
	"deu/internal/repository"
	"deu/internal/users"
	"deu/pkg/db"
	"deu/pkg/router"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	var logLevel slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	if cfg.DatabaseURL == "" {
		log.Fatal("database_url is required in config")
	}
	gormDB := db.InitDB(cfg.DatabaseURL)

	var userRepo repository.UserRepository = repository.NewPostgresUserRepository(gormDB)
	var placeRepo repository.PlaceRepository = repository.NewPostgresPlaceRepository(gormDB)
	var userPlaceRepo repository.UserPlaceRepository = repository.NewPostgresUserPlaceRepository(gormDB)

	if cfg.EnableRequestLogging {
		userRepo = repository.NewLoggingUserRepository(userRepo, logger)
		placeRepo = repository.NewLoggingPlaceRepository(placeRepo, logger)
		userPlaceRepo = repository.NewLoggingUserPlaceRepository(userPlaceRepo, logger)
	}

	userService := users.NewUserService(userRepo, userPlaceRepo, placeRepo)
	placeService := places.NewPlaceService(placeRepo, cfg.EnableCache)

	userHandler := &users.Handler{Service: userService}
	placeHandler := &places.Handler{
		Service:       placeService,
		AllowDeletion: cfg.AllowPlaceDeletion,
	}

	r := router.NewRouter(router.Config{
		UserHandler:  userHandler,
		PlaceHandler: placeHandler,
	})

	serverPort := cfg.ServerPort
	if serverPort == "" {
		serverPort = "8080"
	}

	var handler http.Handler = r
	if cfg.EnableRequestLogging {
		handler = requestLoggingMiddleware(handler)
	}

	if cfg.MaxConnections > 0 {
		handler = maxConnectionsMiddleware(handler, cfg.MaxConnections)
	}

	next := handler
	handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if req.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, req)
	})

	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: handler,
	}
	if cfg.RequestTimeoutSeconds > 0 {
		srv.ReadTimeout = time.Duration(cfg.RequestTimeoutSeconds) * time.Second
		srv.WriteTimeout = time.Duration(cfg.RequestTimeoutSeconds) * time.Second
		srv.IdleTimeout = time.Duration(cfg.RequestTimeoutSeconds) * 2 * time.Second
		slog.Info("Server timeouts configured", "timeout_seconds", cfg.RequestTimeoutSeconds)
	}

	slog.Info("Server starting",
		"port", serverPort,
		"cache_enabled", cfg.EnableCache,
		"max_connections", cfg.MaxConnections,
		"request_logging", cfg.EnableRequestLogging)

	log.Fatal(srv.ListenAndServe())
}

func requestLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)
		slog.Info("HTTP Request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", wrapped.statusCode,
			"duration_ms", duration.Milliseconds(),
			"remote_addr", r.RemoteAddr,
		)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
func maxConnectionsMiddleware(next http.Handler, maxConns int) http.Handler {
	semaphore := make(chan struct{}, maxConns)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case semaphore <- struct{}{}:
			defer func() { <-semaphore }()
			next.ServeHTTP(w, r)
		default:
			http.Error(w, "Server too busy - max connections reached", http.StatusServiceUnavailable)
			slog.Warn("Connection rejected - max connections reached", "max", maxConns)
		}
	})
}