package handlers

import (
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
)

type PostTaskHandler struct {
	taskStore store.TaskStore
}

type PostTaskHandlerParams struct {
	TaskStore store.TaskStore
}

func NewPostTaskHandler(params PostTaskHandlerParams) *PostTaskHandler {
	return &PostTaskHandler{
		taskStore: params.TaskStore,
	}
}

func (h *PostTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	description := r.FormValue("email")

	err := h.taskStore.CreateTask(description)

	if err != nil {

		w.WriteHeader(http.StatusBadRequest)
		// c := templates.TaskError()
		// c.Render(r.Context(), w)
		return
	}

	c := templates.TaskSuccess(description)
	err = c.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}

}

