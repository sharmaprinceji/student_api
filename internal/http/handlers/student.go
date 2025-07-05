package Student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/sharmaprinceji/student-api/internal/types"
	"github.com/sharmaprinceji/student-api/internal/utils/response"
)

func New() http.HandlerFunc{
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




		slog.Info("creaing a student")
		response.WriteJSON(w, http.StatusCreated, map[string]string{"success": "student created successfully"})	
	}
}
