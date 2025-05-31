package main

import (
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
)

func (s *Server) MountRoutes() {

	s.Router.Route("/api/v1", func(r chi.Router) {
		// graphQl
		r.Route("/graphql", func(r chi.Router) {
			r.Handle("/", s.GQlServer)
			r.Handle("/platground", playground.Handler("GraphQL playground", "/api/v1/graphql"))
		})

	})
}
