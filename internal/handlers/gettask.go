package handlers

import (
	"goth/internal/store"
	"goth/internal/templates"
	"goth/internal/middleware"
	"net/http"
	"github.com/go-chi/chi/v5"
	"strconv"
)

type GetTaskHandler struct {
	taskStore store.TaskStore
}

type GetTaskHandlerParams struct {
	TaskStore store.TaskStore
}

func NewGetTaskHandler(params GetTaskHandlerParams) *GetTaskHandler {
	return &GetTaskHandler{
		taskStore: params.TaskStore,
	}
}

func BadRequest(w http.ResponseWriter, r *http.Request) error {
	 	w.WriteHeader(http.StatusBadRequest)
	 	c := templates.TaskError()
	 	c.Render(r.Context(), w)
        return nil
}

func (h *GetTaskHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    u64, err := strconv.ParseUint(id, 10, 32)


    if err != nil { BadRequest(w, r); return}

    sessionValue := r.Context().Value(middleware.SessionKey)
    if sessionValue == nil { BadRequest(w, r); return}

    sessionID := sessionValue.(string)

    ID := uint(u64)
    task, err := h.taskStore.GetTask(ID, sessionID)

    if err != nil { BadRequest(w, r); return}

    taskPage := templates.TaskPage(id, task.Description)
	err = templates.HomeLayout(taskPage, "ToDo App").Render(r.Context(), w)

	if err != nil {
		http.Error(w, "error rendering template", http.StatusInternalServerError)
		return
	}
}
