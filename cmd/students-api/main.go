package main

import (
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	
	"github.com/Naman151/Go-api/internal/config"
	"github.com/Naman151/Go-api/internal/http/handlers/student"
	"github.com/Naman151/Go-api/internal/storage/sqlite"
)

func main() {
	// Load Config
	cfg := config.MustLoad()

	// database setup
	_, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("Storage Initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	//setup router
	router := http.NewServeMux()

	slog.Info("Server Working")
	router.HandleFunc("GET /api/students", student.Create())
	
	//setup server
	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	done := make(chan os.Signal, 1)
	
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	go func() {
		err := server.ListenAndServe()
		slog.Info("Server Working %s", slog.String("address", cfg.Addr))
		if err != nil {
			log.Fatalf("Failed to Start Server %s", err.Error())
		}
	}()
	
	<-done
	
	slog.Info("Shutting Down the Sever")
	server.Shutdown()
}
