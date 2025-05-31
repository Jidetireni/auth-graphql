package main

import (
	"auth-graphql/config"
	"auth-graphql/graph"
	"auth-graphql/repository"
	"log/slog"
	"net/http"
	"os"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/vektah/gqlparser/v2/ast"
)

type Server struct {
	Config    *config.Config
	DB        *repository.Queries
	GQlServer *gqlhandler.Server
	Router    *chi.Mux
	RawDB     sqlx.DB
	Logger    *slog.Logger
}

func NewServer() *Server {
	// Load configuration
	cfg, err := config.New()
	if err != nil {
		panic("Failed to load configuration: " + err.Error())
	}

	// Initialize database connection
	// db, err := repository.DatabaseInit(cfg)
	// if err != nil {
	// 	panic("Failed to initialize database: " + err.Error())
	// }

	// Initialize router
	r := chi.NewRouter()

	// logger setup
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	logger := slog.New(jsonHandler)

	return &Server{
		Config: cfg,
		Router: r,
		Logger: logger,
	}
}

func (s *Server) Start() {
	// Initialize database connection
	db, err := repository.DatabaseInit(s.Config)
	if err != nil {
		s.Logger.Error("Failed to initialize database", "error", err)
		return
	}
	s.DB = repository.New(db)

	// Set up routes
	s.Router.Handle("/", s.GQlServer)

	// Start the server
	s.Logger.Info("Starting server", "host", s.Config.Server.Host, "port", s.Config.Server.Port)
	if err := http.ListenAndServe(s.Config.Server.Host+":"+s.Config.Server.Port, s.Router); err != nil {
		s.Logger.Error("Failed to start server", "error", err)
	}
}

func (s *Server) SetGqlServer() {

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

}
