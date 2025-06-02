package main

import (
	"auth-graphql/config"
	"auth-graphql/graph"
	"auth-graphql/repository"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
	"github.com/vektah/gqlparser/v2/ast"
)

type Server struct {
	Config     *config.Config
	DB         *repository.Queries
	GQlServer  *gqlhandler.Server
	Router     *chi.Mux
	HTTPServer *http.Server
	Logger     *slog.Logger
}

func NewServer() (*Server, error) {
	cfg, err := config.New()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	r := chi.NewRouter()
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(jsonHandler)

	httpServer := &http.Server{
		Addr:         cfg.Server.Host + ":" + cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		Config:     cfg,
		Router:     r,
		HTTPServer: httpServer,
		Logger:     logger,
	}, nil
}

func (s *Server) Initialize() error {
	db, err := repository.DatabaseInit(s.Config)
	if err != nil {
		s.Logger.Error("Failed to initialize database", "error", err)
		return fmt.Errorf("database initialization failed: %w", err)
	}

	s.DB = repository.New(db)
	if err := s.SetGqlServer(); err != nil {
		return fmt.Errorf("GraphQL server setup failed: %w", err)
	}
	s.MountRoutes()

	return nil
}

func (s *Server) Start() error {
	// Create a channel to listen for interrupt signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		s.Logger.Info("Starting server",
			"host", s.Config.Server.Host,
			"port", s.Config.Server.Port,
			"addr", s.HTTPServer.Addr,
		)

		if err := s.HTTPServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.Logger.Error("Failed to start server", "error", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	s.Logger.Info("Shutting down server...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server gracefully
	if err := s.HTTPServer.Shutdown(ctx); err != nil {
		s.Logger.Error("Failed to shutdown server gracefully", "error", err)
		return err
	}

	s.Logger.Info("Server stopped")
	return nil
}

func (s *Server) SetGqlServer() error {
	c := graph.Config{
		Resolvers: &graph.Resolver{},
	}

	srv := gqlhandler.New(graph.NewExecutableSchema(c))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	s.GQlServer = srv
	return nil
}
