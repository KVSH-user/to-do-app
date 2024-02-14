package todo

import (
	"errors"
	"fmt"
	"github.com/KVSH-user/to-do-app/internal/http-server/handlers/auth"
	resp "github.com/KVSH-user/to-do-app/internal/lib/api/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AddTask struct {
	Task string `json:"task"`
}

type CompletedTask struct {
	Id int `json:"id"`
}

type EditTask struct {
	Id         int    `json:"id"`
	EditedTask string `json:"edited_task"`
}

type DelTask struct {
	Id int `json:"id"`
}

type Response struct {
	Id        int           `json:"id"`
	Task      string        `json:"task,omitempty"`
	CreatedAt time.Duration `json:"created_at,omitempty"`
	Active    bool          `json:"active,omitempty"`
	Status    string        `json:"status,omitempty"`
}

type TaskList struct {
	Id        int       `json:"id"`         // Соответствует полю id в таблице
	Task      string    `json:"task"`       // Соответствует полю task в таблице
	Active    bool      `json:"active"`     // Соответствует полю active в таблице
	CreatedAt time.Time `json:"created_at"` // Соответствует полю created_at в таблице
}

type Task interface {
	Create(task string, uid int) (int, error)
	Delete(id, uid int) error
	Complete(id int) error
	Edit(id int, editedTask string) error
	GetAll(uid int) ([]TaskList, error)
}

func GetTasks(log *slog.Logger, task Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.todo.GetTasks"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Errorf("authorization header is missing")
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			fmt.Errorf("invalid token format")
		}

		tokenString := splitToken[1]

		uid, token, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_ = token

		uidInt, err := strconv.Atoi(uid)

		tasks, err := task.GetAll(uidInt)
		if err != nil {
			log.Error("failed to get tasks: ", err)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		render.JSON(w, r, tasks)
	}
}

func Create(log *slog.Logger, task Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.todo.Create"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Errorf("authorization header is missing")
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			fmt.Errorf("invalid token format")
		}

		tokenString := splitToken[1]

		uid, token, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_ = token

		uidInt, err := strconv.Atoi(uid)

		var req AddTask

		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		id, err := task.Create(req.Task, uidInt)
		if err != nil {
			log.Error("failed to create task: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("task created")

		render.JSON(w, r, Response{
			Id:   id,
			Task: req.Task,
		})

	}
}

func Delete(log *slog.Logger, task Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.todo.Delete"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			fmt.Errorf("authorization header is missing")
		}

		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			fmt.Errorf("invalid token format")
		}

		tokenString := splitToken[1]

		uid, token, err := auth.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		_ = token

		uidInt, err := strconv.Atoi(uid)

		var req DelTask

		err = render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err = task.Delete(req.Id, uidInt)
		if err != nil {
			log.Error("failed to delete task: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("task deleted")

		render.JSON(w, r, Response{
			Id:     req.Id,
			Status: "task deleted",
		})

	}
}

func Complete(log *slog.Logger, task Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.todo.Complete"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req CompletedTask

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err = task.Complete(req.Id)
		if err != nil {
			log.Error("failed to complete task: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("task completed")

		render.JSON(w, r, Response{
			Id:     req.Id,
			Status: "task completed",
		})

	}
}

func Edit(log *slog.Logger, task Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.todo.Edit"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req EditTask

		err := render.DecodeJSON(r.Body, &req)
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")

			render.JSON(w, r, resp.Error("empty request"))

			return
		}
		if err != nil {
			log.Error("failed to decode request body: ", err)

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		err = task.Edit(req.Id, req.EditedTask)
		if err != nil {
			log.Error("failed to edit task: ", err)

			render.JSON(w, r, resp.Error("internal error"))

			return
		}

		log.Info("task edited")

		render.JSON(w, r, Response{
			Id:     req.Id,
			Task:   req.EditedTask,
			Status: "task edited",
		})

	}
}
