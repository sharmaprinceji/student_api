package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sharmaprinceji/student-api/internal/config"
	"github.com/sharmaprinceji/student-api/internal/http/handlers"
	"github.com/sharmaprinceji/student-api/internal/storage/sqlite"
)

func main(){
	// fmt.Println("Hello, Student API!")
	//loadConfig()
	cfg := config.MustLoad()
	// database setup
	storage,er:=sqlite.New(cfg)

	if er!=nil {
		log.Fatal(er)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version","1.0.0"))

	//setup router
	router:=http.NewServeMux()

	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    //   w.Write([]byte("Welcome to Student server!"))
	// })

	router.HandleFunc("POST /api/student",Student.New(storage))
	router.HandleFunc("GET /api/student/{id}",Student.GetById(storage))
	router.HandleFunc("GET /api/students", Student.GetAll(storage)) // Get all students
	router.HandleFunc("GET /api/students/{page}", Student.GetListPagination(storage)) 
	router.HandleFunc("PUT /api/student/{id}", Student.UpdateById(storage))
	router.HandleFunc("DELETE /api/student/{id}", Student.DeleteById(storage))

	//setup server.
	server:= http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      router,
	}

	slog.Info("Starting server...", slog.String("address", cfg.HTTPServer.Addr))
   
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