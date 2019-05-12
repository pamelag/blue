package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pamelag/blue/article"
	"github.com/pamelag/blue/content"
	"github.com/pamelag/blue/postgres"
	"github.com/pamelag/blue/server"
	"github.com/pamelag/blue/tag"
)

func main() {

	var (
		as article.Service
		ts tag.Service

		articles content.ArticleRepository
		tags content.TagRepository
	)

	// create connection pool
	pool, err := getConnPool()
	if err != nil {
		panic(err)
	}

	log.Println("Created connection pool successfully")

	// Setup repository
	articles = postgres.NewArticleRepository(pool)
	tags = postgres.NewTagRepository(pool)

	// Inject repository into services
	as = article.NewService(articles)
	ts = tag.NewService(tags)

	// Inject services to routing server 
	srv := server.New(as, ts)
	httpAddr := ":" + port

	httpServer := &http.Server{Addr: httpAddr,
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second}

	// registers a function to call on Shutdown
	httpServer.RegisterOnShutdown(func() {
		log.Println("Call shutdown hooks")
	})

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGINT)
		signal.Notify(sigint, syscall.SIGTERM)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := httpServer.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed

}
