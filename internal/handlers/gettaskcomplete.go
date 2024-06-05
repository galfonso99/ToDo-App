package handlers

import (
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"github.com/go-chi/chi/v5"
	"strconv"
)

type GetTaskCompleteHandler struct {
	taskStore store.TaskStore
}

type GetTaskCompleteHandlerParams struct {
	TaskStore store.TaskStore
}

func NewGetTaskCompleteHandler(params GetTaskCompleteHandlerParams) *GetTaskCompleteHandler {
	return &GetTaskCompleteHandler{
		taskStore: params.TaskStore,
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) error {
	 	w.WriteHeader(http.StatusBadRequest)
	 	c := templates.TaskError()
	 	c.Render(r.Context(), w)
        return nil
}

func (h *GetTaskCompleteHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    u64, err := strconv.ParseUint(id, 10, 32)

    if err != nil { BadRequest(w, r); return}

    ID := uint(u64)
    task, err := h.taskStore.GetTask(ID)

    if err != nil { BadRequest(w, r); return}

    completedTask := templates.CompletedTodo(id, task.Description)
	err = completedTask.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}