package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/alphacoder0x/go-rest-api/internal/storage"
	"github.com/alphacoder0x/go-rest-api/internal/types"
	"github.com/alphacoder0x/go-rest-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

type Student struct {
	name string
}

func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("creating a student")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)

		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("Empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		lastid, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, err.Error())
			return
		}

		response.WriteJson(w, http.StatusCreated, map[string]any{"status": "ok", "generated_id": lastid})
	}
}

func GetStudents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("returning students"))
	}
}
