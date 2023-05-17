package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/render"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/rifqoi/aws-project/todo-backend/common"
)

type TodoHandlers struct {
	todoSvc TodoService
}

func NewTodoHandlers(todoSvc TodoService) *TodoHandlers {
	if todoSvc == nil {
		panic("todosvc is nil")
	}

	return &TodoHandlers{todoSvc}
}

type TaskCreateRequest struct {
	Task string `json:"task,omitempty"`
}

// Convenient way from go-chi to validate request body
func (t TaskCreateRequest) Bind(r *http.Request) error {
	task := validation.Field(&t.Task, validation.Required)
	err := validation.ValidateStruct(&t, task)
	if err != nil {
		return common.WrapErrorf(err, common.ErrorTypeInvalidInput, "invalid-request")
	}

	return nil
}

func (t *TodoHandlers) AddTask(w http.ResponseWriter, r *http.Request) {

	var req TaskCreateRequest
	err := render.Bind(r, &req)
	if err != nil {
		RespondHTTPError(err, "invalid request", w, r)
		return
	}

	todo, err := t.todoSvc.AddTask(r.Context(), req.Task)
	if err != nil {
		RespondHTTPError(err, "failed to add task", w, r)
		return
	}

	resp := NewSuccessResponse("successfully create a task", todo)
	render.JSON(w, r, resp)

}

// This is function to decode and encode json from request and response
// but I ain't gonna use it because DRY ;p
func decodeJson(r *http.Request, data any) error {
	err := json.NewDecoder(r.Body).Decode(data)
	return err
}

func encodeJSON(w http.ResponseWriter, v any) error {
	err := json.NewEncoder(w).Encode(v)
	return err
}
