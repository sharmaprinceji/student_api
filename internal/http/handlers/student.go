package Student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	 "strconv"

	"github.com/go-playground/validator/v10"
	"github.com/sharmaprinceji/student-api/internal/storage"
	"github.com/sharmaprinceji/student-api/internal/types"
	"github.com/sharmaprinceji/student-api/internal/utils/response"
)

func New(storage storage.Storage)  http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {

		var student types.Student
		err:=json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			// response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "request body is empty"})
			response.WriteJSON(w, http.StatusBadRequest,response.GeneralError(fmt.Errorf("request body is Invalid")))
			return
		}

		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to decode request body: %v", err)))
			return 
		}


		if err:=validator.New().Struct(student); err != nil {
			validateErrs:= err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return 
		}

        id,err:=storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
			student.City,
		)

		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError,err)
			return
		}


		slog.Info("student created successfully with-",slog.Int64("id", id))
		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})	
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}
        intid,err:=strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format: %v", err)))
			return
		}

		student, err := storage.GetStudentById(intid)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get student: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusOK, student)
	}
}
