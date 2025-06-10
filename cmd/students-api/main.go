package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Naman151/Go-api/internal/config"
	"github.com/Naman151/Go-api/internal/http/handlers/student"
	"github.com/Naman151/Go-api/internal/storage/sqlite"
)

func gracefulShutdown(server *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	done <- true
}

func main() {
	// Load Config
	cfg := config.MustLoad()

	// database setup
	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router
	router := http.NewServeMux()

	// slog.Info("Server Working", slog.String("address", cfg.Addr))
	router.HandleFunc("POST /api/students", student.Create(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students/", student.GetList(storage))
	router.HandleFunc("DELETE /api/students/{id}", student.DeleteById(storage))
	// router.HandleFunc("PATCH /api/students/{id}", student.UpdateById(storage))

	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan bool, 1)

	// signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go gracefulShutdown(&server, done)

	err = server.ListenAndServe()
	slog.Info("Server Working %s", slog.String("address", cfg.Addr))
	if err != nil {
		log.Fatalf("Failed to Start Server %s", err.Error())
	}

	<-done
	// slog.Info("Shutting Down the Sever")

	slog.Info("server shutdown successfully")
}
