package studentrouter

import (
	"log"
	"net/http"
	"github.com/sharmaprinceji/student-api/internal/router"
	"github.com/sharmaprinceji/student-api/internal/http/handlers"
	"github.com/sharmaprinceji/student-api/db"
	"github.com/sharmaprinceji/student-api/internal/config"
)


func StudentRouter() *http.ServeMux {
	cfg := config.MustLoad()

	storage, err := db.Mydb(cfg)
	if err != nil {
		log.Fatalf("failed to init DB: %v", err)
	}

	route :=  router.StudentRoute()

	route.HandleFunc("GET /api/student/{id}", Student.GetById(storage))
	route.HandleFunc("GET /api/students", Student.GetAll(storage)) // Get all students
	route.HandleFunc("GET /api/students/{page}", Student.GetListPagination(storage))
	route.HandleFunc("PUT /api/student/{id}", Student.UpdateById(storage))
	route.HandleFunc("DELETE /api/student/{id}", Student.DeleteById(storage))

	return route;
}