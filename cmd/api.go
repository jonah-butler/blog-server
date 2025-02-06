package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"blog-api/db"
	handlers "blog-api/handlers/blog"
	repositories "blog-api/repositories/blog"
	services "blog-api/services/blog"

	"github.com/joho/godotenv"
)

func main() {

	// load env
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("Unable to load env: ", err.Error())
	}

	// connect to db
	uri, hasURI := os.LookupEnv("MONGO_DB_URI")
	if !hasURI {
		log.Fatal("Unable to load database URI")
	}

	dbName, hasDbName := os.LookupEnv("MONGO_DB_NAME")
	if !hasDbName {
		log.Fatal("Unable to load database name")
	}

	db, err := db.ConnecToMongo(uri, dbName)
	if err != nil {
		log.Fatal("Unable to connect to database: " + err.Error())
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := db.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from database: %v", err)
		}
	}()

	log.Println("connected to database")

	port, hasPort := os.LookupEnv("PORT")
	if !hasPort {
		log.Fatal("Port value unavailable")
	}

	// initialize repos
	blogRepo := repositories.NewBlogRepository(db.DB)

	// initialize services
	blogService := services.NewBlogService(blogRepo)

	// initialize handlers
	blogController := handlers.NewBlogHandler(blogService)

	// initialize server
	mux := http.NewServeMux()

	blogController.RegisterBlogRoutes("/blog", mux)

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// run server inside channel
	go func() {
		fmt.Printf("listening on port %s", port)
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down...")
}
