package main

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
)

func (s *Server) MountRoutes() {
	s.Router.Route("/api/v1", func(r chi.Router) {
		r.Route("/graphql", func(r chi.Router) {
			r.Handle("/", s.GQlServer)
			r.Handle("/playground", playground.Handler("GraphQL Playground", "/api/v1/graphql"))
		})
	})
}
