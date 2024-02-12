package todo

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"time"
	resp "to-do-app/internal/lib/api/response"
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
	Create(task string) (int, error)
	Delete(id int) error
	Complete(id int) error
	Edit(id int, editedTask string) error
	GetAll() ([]TaskList, error)
}

func GetTasks(log *slog.Logger, task Task) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.todo.GetTasks"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		tasks, err := task.GetAll()
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

		var req AddTask

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

		id, err := task.Create(req.Task)
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

		var req DelTask

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

		err = task.Delete(req.Id)
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
