package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/erknas/forum/graph"
	"github.com/erknas/forum/internal/config"
	"github.com/erknas/forum/internal/service"
	"github.com/erknas/forum/internal/storage"
	"github.com/erknas/forum/internal/subscription"
	"github.com/erknas/forum/migrations/migrator"
)

func main() {
	var (
		ctx      = context.Background()
		cfg      = config.Load()
		svc      *service.Service
		sub      = subscription.New()
		inmemory = os.Getenv("IN_MEMORY")
	)

	if inmemory == "false" {
		if err := migrator.New(cfg); err != nil {
			log.Fatalf("failed to migrate: %s", err)
		}

		postgres, err := storage.NewPostgresPool(ctx, cfg)
		if err != nil {
			log.Fatalf("failed to connect to postgres: %s", err)
		}

		svc = service.New(postgres, sub)

		log.Println("using postgres storage")
	} else {
		inmemory := storage.NewInMemoryStorage()
		svc = service.New(inmemory, sub)

		log.Println("using in-memory storage")
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{Svc: svc, Sub: sub}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{})

	srv.Use(extension.Introspection{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("starting server on [http://localhost:%s/]", cfg.Addr)
	log.Fatal(http.ListenAndServe(":"+cfg.Addr, nil))
}
