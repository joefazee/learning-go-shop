package server

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/joefazee/learning-go-shop/graph"
	"github.com/joefazee/learning-go-shop/graph/resolver"
	"github.com/vektah/gqlparser/v2/ast"
)

func (s *Server) createGraphQLHandler() *handler.Server {

	rvr := resolver.NewResolver(
		s.authService,
		s.userService,
		s.productService, s.cartService,
		s.orderService,
	)

	schema := graph.NewExecutableSchema(graph.Config{Resolvers: rvr})

	srv := handler.New(schema)

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return srv
}

func (s *Server) graphqlHandler() gin.HandlerFunc {
	h := s.createGraphQLHandler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *Server) playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL Playground", "/graphql/")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
