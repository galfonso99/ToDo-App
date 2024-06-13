package handlers

import (
	"goth/internal/middleware"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"strconv"
	"time"
	"github.com/go-chi/chi/v5"
)

type PutTaskHandler struct {
	taskStore store.TaskStore
}

type PutTaskHandlerParams struct {
	TaskStore store.TaskStore
}

func NewPutTaskHandler(params PutTaskHandlerParams) *PutTaskHandler {
	return &PutTaskHandler{
		taskStore: params.TaskStore,
	}
}

func (h *PutTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    time.Sleep(500 * time.Millisecond)
    description := r.FormValue("task")
    idString := chi.URLParam(r, "id")
    idU32, err := strconv.ParseUint(idString, 10, 32)

    if err != nil { BadRequest(w, r); return}

    id := uint(idU32)

    sessionValue := r.Context().Value(middleware.SessionKey)
    if sessionValue == nil { BadRequest(w, r); return}

    sessionID := sessionValue.(string)

    err = h.taskStore.EditTask(id, sessionID, description)

    if err != nil { BadRequest(w, r); return}

    taskEditor := templates.TaskDescription(idString, description)
	err = taskEditor.Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	} 
}
