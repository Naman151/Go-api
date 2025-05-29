package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Naman151/Go-api/internal/config"
)

func main()  {
	// Load Config
	cfg := config.MustLoad()

	//setup router
	router := http.NewServeMux()

	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Server Working")
		w.Write([]byte("Welcome to Students Api"))
	})

	//setup server
	server := http.Server{
		Addr: cfg.Addr,
		Handler: router,
	}

	err := server.ListenAndServe()

	if err != nil {
		log.Fatalf("Failed to Start Server %s", err.Error())
	}
	defer server.Close()
}
