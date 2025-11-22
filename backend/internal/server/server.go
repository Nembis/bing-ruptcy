package server

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/nembis/bing-ruptcy/backend/internal/config"
	"github.com/nembis/bing-ruptcy/backend/internal/database"
	"github.com/rs/cors"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var name = "github.com/nembis/bing-ruptcy/backend"

var baseLogger = otelslog.NewLogger(name)

type server struct {
	queries *database.Queries
	cfg     *config.Config
}

func NewServer(db database.DBTX, cfg *config.Config) *server {
	return &server{
		queries: database.New(db),
		cfg:     cfg,
	}
}

func (s *server) GetHttpServelr(ctx context.Context) *http.Server {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	router := http.NewServeMux()
	handleFunc := func(pattern string, handlerFunc func(http.ResponseWriter, *http.Request)) {
		coreHandler := http.HandlerFunc(handlerFunc)
		tracedHandler := otelhttp.NewHandler(coreHandler, pattern)
		handler := otelhttp.WithRouteTag(pattern, tracedHandler)
		router.Handle(pattern, handler)
	}

	handleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		logger := getEnrichedLogger(r.Context())
		logger.InfoContext(r.Context(), "Checking health")
		respondWithJSON(w, http.StatusOK, struct{}{})
	})

	handler := c.Handler(router)
	handler = midddlewareLoggerEnricher(handler)

	serve := http.Server{
		Addr:         ":4000",
		Handler:      handler,
		BaseContext:  func(_ net.Listener) context.Context { return ctx },
		ReadTimeout:  time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &serve
}
