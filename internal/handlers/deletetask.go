package handlers

import (
	"fmt"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type DeleteTaskHandler struct {
	taskStore store.TaskStore
}

type DeleteTaskHandlerParams struct {
	TaskStore store.TaskStore
}

func NewDeleteTaskHandler(params DeleteTaskHandlerParams) *DeleteTaskHandler {
	return &DeleteTaskHandler{
		taskStore: params.TaskStore,
	}
}

func (h *DeleteTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

    id := chi.URLParam(r, "id")
    fmt.Println(id)
    u64, err := strconv.ParseUint(id, 10, 32)

    if err != nil {
	 	w.WriteHeader(http.StatusBadRequest)
	 	c := templates.TaskError()
	 	c.Render(r.Context(), w)
	 	return
    }

    ID := uint(u64)
    err = h.taskStore.DeleteTask(ID)

	 if err != nil {
	 	w.WriteHeader(http.StatusBadRequest)
	 	c := templates.TaskError()
	 	c.Render(r.Context(), w)
	 	return
	 }
    emptyTask := templates.DeleteTask()
	err = emptyTask.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}


