package graphql

import (
	"context"
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
	"github.com/markbates/goth/gothic"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func init() {
	initializers.ConnectToDatabase()
}

// Must contain 6 characters, one uppercase, one lowercase, one number, and one special character
func Handler() {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},  // Allow only your frontend origin
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"}, // Allow common methods
		AllowedHeaders:   []string{"*"},                      // Allow all headers
		ExposedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
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
	mux := http.NewServeMux()

	// GraphQL Playground & API
	mux.Handle("/playground", c.Handler(auth.Middleware(playground.Handler("GraphQL playground", "/"))))
	mux.Handle("/graphql", c.Handler(auth.Middleware(srv)))

	// File Upload & Download
	mux.HandleFunc("/upload", auth.MiddlewareFuncForUploads(uploadFileHandler))
	mux.HandleFunc("/download", auth.MiddlewareFuncForUploads(downloadFileHandler))

	// Static File Serving
	mux.Handle("/", http.FileServer(http.Dir("static")))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Document Listing API
	mux.HandleFunc("/documents", listDocuments)

	// Serve static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Serve HTML pages
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Serving /hello")
		http.ServeFile(w, r, "./static/hello.html")
	})

	// Redirect root to /hello
	// mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "/hello", http.StatusSeeOther)
	// })

	mux.HandleFunc("/oauth-success", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/0Auth.html")
	})
	// mux.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
	// 	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/auth/"), "/")
	// 	if len(pathParts) > 0 && pathParts[0] != "" {
	// 		gothic.BeginAuthHandler(w, r)
	// 		return
	// 	}
	// 	http.Error(w, "Invalid OAuth request", http.StatusNotFound)
	// })

	mux.HandleFunc("/auth/{provider}/callback", auth.OauthCallbackHandler)
	mux.HandleFunc("/auth/", func(w http.ResponseWriter, r *http.Request) {
		provider := r.URL.Path[len("/auth/"):]
		if provider == "" {
			http.Error(w, "You must select a provider", http.StatusBadRequest)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "provider", provider))
		gothic.BeginAuthHandler(w, r)
	})

	// Start server
	log.Println("Server running on http://localhost:8080")

	// Start Server
	log.Printf("GraphQL running at http://localhost:%s/graphql", port)
	log.Printf("File Upload running at http://localhost:%s/upload", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
