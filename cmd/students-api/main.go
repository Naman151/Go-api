package main

import (
	"log"
	"log/slog"
	"net/http"

	"github.com/Naman151/Go-api/internal/config"
)

func main()  {
	// Load Config
	cfg := config.MustLoad()

	//setup router
	router := http.NewServeMux()

	slog.Info("Server Working")
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to Students Api"))
	})

	//setup server
	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	err := server.ListenAndServe()
	slog.Info("Server Working %s", slog.String("address", cfg.Addr))
	if err != nil {
		log.Fatalf("Failed to Start Server %s", err.Error())
	}
	defer server.Close()
}
