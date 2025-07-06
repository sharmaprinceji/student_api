package main

import (
	"context"
	"fmt"

	//"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/sharmaprinceji/student-api/db"
	"github.com/sharmaprinceji/student-api/internal/config"
	studentrouter "github.com/sharmaprinceji/student-api/internal/router/studentRouter"
	// "github.com/sharmaprinceji/student-api/internal/http/handlers"
	"github.com/sharmaprinceji/student-api/internal/http/schedular"
	// "github.com/sharmaprinceji/student-api/internal/router"
	//"github.com/sharmaprinceji/student-api/internal/router/studentrouter"
	// "github.com/sharmaprinceji/student-api/internal/storage/sqlite"
)

func main() {
	// fmt.Println("Hello, Student API!")
	//loadConfig()
	cfg := config.MustLoad() /// 👈 call only once here

	// database setup
	// storage,er:=sqlite.New(cfg)
	// if er!=nil {
	// 	log.Fatal(er)
	// }
	// slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version","1.0.0"))

	//seprate db setup ..
	// _, err := db.Mydb(cfg)
	// if err != nil {
	// 	log.Fatalf("failed to init DB: %v", err)
	// }

	// log.Println("Db connection on..", cfg.HTTPServer.Addr)

	//setup router
	//router:=http.NewServeMux()
	//route :=  router.StudentRoute()
	route:=studentrouter.StudentRouter();

	// router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//   w.Write([]byte("Welcome to Student server!"))
	// route.HandleFunc("POST /api/student", Student.New(storage))


	// route.HandleFunc("GET /api/student/{id}", Student.GetById(storage))
	// route.HandleFunc("GET /api/students", Student.GetAll(storage)) // Get all students
	// route.HandleFunc("GET /api/students/{page}", Student.GetListPagination(storage))
	// route.HandleFunc("PUT /api/student/{id}", Student.UpdateById(storage))
	// route.HandleFunc("DELETE /api/student/{id}", Student.DeleteById(storage))

	//setup server.
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: route,
	}

	//StartCronJob
	//scheduler.StartCronJob()
	scheduler.StartStudentFetchJob()

	slog.Info("Starting server...", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			return
		}
	}()

	<-done // Wait for a signal to stop the server

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	er := server.Shutdown(ctx)
	if er != nil {
		slog.Error("failed to shutting down server", slog.String("error", er.Error()))
	}

	slog.Info("Server stopped gracefully")

}
