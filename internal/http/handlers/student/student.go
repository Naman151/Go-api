package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Naman151/Go-api/internal/storage"
	"github.com/Naman151/Go-api/internal/types"
	"github.com/Naman151/Go-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func Create(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Creating New Student")
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, err.Error())
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		//---------REQUEST VALIDATION --------------------------
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
		}

		lastId, err := storage.CreateStudent(student.Name, student.Email, student.Age)

		if err != nil {
			slog.Info("error ", slog.String("error", fmt.Sprint(err.Error())))
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info("Student Created Sucessfully", slog.String("userId", fmt.Sprint(lastId)))

		response.WriteJson(w, http.StatusCreated, map[string]int64{"Id": lastId})
	}
}

func GetById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Student Id", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		student, err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, student)
	}

}

func GetList(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		students, err := storage.GetStudentsList()
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, students)
	}

}

// func UpdateById(storage storage.Storage) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		id := r.PathValue("id")
// 		slog.Info("Student Id", slog.String("id", id))

// 		intId, err := strconv.ParseInt(id, 10, 64)
// 		if err != nil {
// 			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
// 		}

// 		students, err := storage.UpdateStudentById(intId)
// 		if err != nil {
// 			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}

// 		response.WriteJson(w, http.StatusOK, students)
// 	}

// }

func DeleteById(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		slog.Info("Student Id to delete", slog.String("id", id))

		intId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		err = storage.DeleteStudentById(intId)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		response.WriteJson(w, http.StatusOK, "Deletion Sucessful")
	}

}
