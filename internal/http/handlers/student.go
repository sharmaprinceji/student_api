package Student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/sharmaprinceji/student-api/internal/storage"
	"github.com/sharmaprinceji/student-api/internal/types"
	"github.com/sharmaprinceji/student-api/internal/utils/response"
)

// func New(storage storage.Storage)  http.HandlerFunc{
// 	return func(w http.ResponseWriter, r *http.Request) {

// 		var student types.Student
// 		err:=json.NewDecoder(r.Body).Decode(&student)
// 		if errors.Is(err, io.EOF) {
// 			// response.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "request body is empty"})
// 			response.WriteJSON(w, http.StatusBadRequest,response.GeneralError(fmt.Errorf("request body is Invalid")))
// 			return
// 		}

// 		if err != nil {
// 			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to decode request body: %v", err)))
// 			return
// 		}

// 		if err:=validator.New().Struct(student); err != nil {
// 			validateErrs:= err.(validator.ValidationErrors)
// 			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
// 			return
// 		}

//         id,err:=storage.CreateStudent(
// 			student.Name,
// 			student.Email,
// 			student.Age,
// 			student.City,
// 		)

// 		if err != nil {
// 			response.WriteJSON(w, http.StatusInternalServerError,err)
// 			return
// 		}

// 		slog.Info("student created successfully with-",slog.Int64("id", id))
// 		response.WriteJSON(w, http.StatusCreated, map[string]int64{"id": id})
// 	}
// }


func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student
		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is invalid")))
			return
		}
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to decode request body: %v", err)))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		id, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
			student.City,
		)

		if err != nil {
			// Handle duplicate email error
			if strings.Contains(err.Error(), "UNIQUE constraint failed") || strings.Contains(err.Error(), "duplicate key") {
				response.WriteJSON(w, http.StatusConflict, response.GeneralError(fmt.Errorf("email already exists")))
				return
			}
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create student: %v", err)))
			return
		}

		slog.Info("student created successfully", slog.Int64("id", id))
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

func GetAll(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetAllStudents()
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get students: %v", err)))
			return
		}
		response.WriteJSON(w, http.StatusOK, students)
	}
}


func GetListPagination(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageStr := r.PathValue("page")
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid page number")))
			return
		}

		const limit = 5
		students, totalCount, err := storage.GetStudentsPaginated(page, limit)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get students: %v", err)))
			return
		}

		totalPages := int(math.Ceil(float64(totalCount) / float64(limit)))

		response.WriteJSON(w, http.StatusOK, map[string]interface{}{
			"currentPage": page,
			"totalPages":  totalPages,
			"data":        students,
		})
	}
}

func UpdateById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}
		intid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format: %v", err)))
			return
		}

		var student types.Student
		err = json.NewDecoder(r.Body).Decode(&student)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("failed to decode request body: %v", err)))
			return
		}

		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJSON(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		slog.Info("updating student with id", slog.Int64("id", intid))
		_, err = storage.UpdateStudentById(intid, student.Name, student.Email, student.Age, student.City)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to update student: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "student updated successfully"})
	}
}

func DeleteById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		if id == "" {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("id is required")))
			return
		}
		intid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJSON(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id format: %v", err)))
			return
		}

		slog.Info("deleting student with id", slog.Int64("id", intid))
		 _, err = storage.DeleteStudent(intid)
		if err != nil {
			response.WriteJSON(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to delete student: %v", err)))
			return
		}

		response.WriteJSON(w, http.StatusOK, map[string]string{"message": "student deleted successfully"})
	}
}