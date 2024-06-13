package handlers

import (
	"goth/internal/store"
	"goth/internal/templates"
	"goth/internal/middleware"
	"net/http"
	"github.com/go-chi/chi/v5"
	"strconv"
)

type GetTaskEditorHandler struct {
	taskStore store.TaskStore
}

type GetTaskEditorHandlerParams struct {
	TaskStore store.TaskStore
}

func NewGetTaskEditorHandler(params GetTaskEditorHandlerParams) *GetTaskEditorHandler {
	return &GetTaskEditorHandler{
		taskStore: params.TaskStore,
	}
}

func (h *GetTaskEditorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    idString := chi.URLParam(r, "id")
    idUsize, err := strconv.ParseUint(idString, 10, 32)

    if err != nil { BadRequest(w, r); return}

    id := uint(idUsize)

    sessionValue := r.Context().Value(middleware.SessionKey)
    if sessionValue == nil { BadRequest(w, r); return}

    sessionID := sessionValue.(string)

    task, err := h.taskStore.GetTask(id, sessionID)

    if err != nil { BadRequest(w, r); return}

    taskEditor := templates.TaskEditor(idString, task.Description)
	err = taskEditor.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
