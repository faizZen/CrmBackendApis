package graphql

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	initializers "github.com/Zenithive/it-crm-backend/Initializers"
	"github.com/Zenithive/it-crm-backend/auth"
	"github.com/Zenithive/it-crm-backend/internal/graphql/generated"
	"github.com/Zenithive/it-crm-backend/internal/graphql/schema"
	"github.com/go-chi/cors"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func init() {
	initializers.ConnectToDatabase()
}

// Must contain 6 characters, one uppercase, one lowercase, one number, and one special character
func Handler() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                      // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Allow common methods
		AllowedHeaders:   []string{"*"},                      // Allow all headers
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: false,
	})

	// handler := c.Handler(http.DefaultServeMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &schema.Resolver{}}))
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})
	http.Handle("/", auth.Middleware(playground.Handler("GraphQL playground", "/")))
	http.Handle("/graphql", c.Handler(auth.Middleware(srv)))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
