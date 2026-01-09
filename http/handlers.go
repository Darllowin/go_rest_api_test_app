package http

import (
	"encoding/json"
	"errors"
	"go-todoList/todo"
	"net/http"
	"time"
)

type HTTPHandlers struct {
	todolist *todo.List
}

func NewHTTPHandlers(todolist *todo.List) *HTTPHandlers {
	return &HTTPHandlers{
		todolist: todolist,
	}
}

func (h *HTTPHandlers) sendError(w http.ResponseWriter, message string, code int) {
	errDTO := ErrorDTO{
		Message: message,
		Time:    time.Now(),
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(errDTO.ToString()))
}

/*
pattern : /task
method: POST
info: JSON in HTTP requesr body

succeed:

	status code: 201
	response body: JSON represent created task

failed:

	status code: 400, 409, 500, ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	var taskDTO TaskDTO
	if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := taskDTO.ValidateForCreate(); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	todotask := todo.NewTask(taskDTO.Title, taskDTO.Description)
	if err := h.todolist.AddTask(todotask); err != nil {
		if errors.Is(err, todo.ErrTaskAlredyExists) {
			h.sendError(w, err.Error(), http.StatusConflict)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(todotask, "", "	")
	if err != nil {
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

/*
pattern : /task/{title}
method: GET
info: pattern

succeed:

	status code: 200
	response body: JSON represent found task

failed:

	status code: 400, 404, 500, ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetTask(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")

	task, err := h.todolist.GetTask(title)
	if err != nil {
		if errors.Is(err, todo.ErrTaskNotFound) {
			h.sendError(w, err.Error(), http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(task, "", "	")
	if err != nil {
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

/*
pattern : /tasks
method: GET
info: -

succeed:

	status code: 200
	response body: JSON represent found tasks

failed:

	status code: 400, 500, ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks := h.todolist.ListTask()
	b, err := json.MarshalIndent(tasks, "", "	")
	if err != nil {
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

/*
pattern : /tasks?compleated=true
method: GET
info: query params

succeed:

	status code: 200
	response body: JSON represent found tasks

failed:

	status code: 400, 500, ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleGetAllUncompletedTasks(w http.ResponseWriter, r *http.Request) {
	uncompletedTasks := h.todolist.ListUnCompletedTasks()
	b, err := json.MarshalIndent(uncompletedTasks, "", "	")
	if err != nil {
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

/*
pattern : /tasks/{title}
method: PATCH
info: pattern + JSON in request body

succeed:

	status code: 200
	response body: JSON represent changed task

failed:

	status code: 400, 409, 500, ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleCompletedTask(w http.ResponseWriter, r *http.Request) {
	var completeDTO CompleteTaskDTO
	if err := json.NewDecoder(r.Body).Decode(&completeDTO); err != nil {
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := r.PathValue("title")

	var (
		changedTask todo.Task
		err         error
	)

	if completeDTO.Complete {
		changedTask, err = h.todolist.CompleteTask(title)
	} else {
		changedTask, err = h.todolist.UnCompleteTask(title)
	}

	if err != nil {
		if errors.Is(err, todo.ErrTaskNotFound) {
			h.sendError(w, err.Error(), http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	b, err := json.MarshalIndent(changedTask, "", "	")
	if err != nil {
		h.sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

/*
pattern : /tasks/{title}
method: DELETE
info: pattern

succeed:

	status code: 204
	response body: -

failed:

	status code: 400, 404, 500, ...
	response body: JSON with error + time
*/
func (h *HTTPHandlers) HandleDeleteTask(w http.ResponseWriter, r *http.Request) {
	title := r.PathValue("title")

	if err := h.todolist.DeleteTask(title); err != nil {
		if errors.Is(err, todo.ErrTaskNotFound) {
			h.sendError(w, err.Error(), http.StatusNotFound)
		} else {
			h.sendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
