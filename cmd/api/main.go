package main

import (
	"log"
	"net/http"
	"time"

	"github.com/priyanshu-kun/project-metagross/internal/config"
	"github.com/priyanshu-kun/project-metagross/internal/db"
	"github.com/priyanshu-kun/project-metagross/internal/handlers"
)

func main() {

	cfg := config.MustLoad()

	_, err := db.Connect(cfg.DB_URL)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)

	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", handlers.Healthz)

	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("Database connected!")
	log.Printf("Server is listening on %s", srv.Addr)

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
