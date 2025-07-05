package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sharmaprinceji/student-api/internal/config"
)

func main(){
	fmt.Println("Hello, Student API!")
	//loadConfig()
	cfg := config.MustLoad()
	// database setup

	//setup router
	router:=http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("Welcome to Student API!"))
	})

	//setup server..
	server:= http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      router,
	}

	slog.Info("Starting server...", slog.String("address", cfg.HTTPServer.Addr))
    // fmt.Printf("Server is running on %s\n", cfg.HTTPServer.Addr)
	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT);

	go func() {
		err:=server.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			return
		}
	}()

	<- done // Wait for a signal to stop the server

	slog.Info("Shutting down server...")

	ctx,cancel:=context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
    
	err:=server.Shutdown(ctx)
	if err != nil {
		slog.Error("failed to shutting down server", slog.String("error", err.Error()))
	}

	slog.Info("Server stopped gracefully")

}