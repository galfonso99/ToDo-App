package handlers

import (
	// "goth/internal/middleware"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
)

type HomeHandler struct{
    taskStore store.TaskStore 
}

type HomeHandlerParams struct {
	TaskStore store.TaskStore
}

func NewHomeHandler(params HomeHandlerParams) *HomeHandler {
	return &HomeHandler{
		taskStore: params.TaskStore,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    tasks , _:= h.taskStore.GetAllTasks()
	c := templates.Home(tasks)
	err := templates.HomeLayout(c, "ToDo App").Render(r.Context(), w)


	if err != nil {
		http.Error(w, "Error rendering template", http.StatusInternalServerError)
		return
	}
}
