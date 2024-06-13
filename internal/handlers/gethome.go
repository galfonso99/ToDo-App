package handlers

import (
	// "goth/internal/middleware"
	b64 "encoding/base64"
	// "fmt"
	"goth/internal/middleware"
	"goth/internal/store"
	"goth/internal/templates"
	"net/http"
	"time"
)

type HomeHandler struct{
    taskStore store.TaskStore 
    sessionStore store.SessionStore
}

type HomeHandlerParams struct {
	TaskStore store.TaskStore
    SessionStore store.SessionStore
}

func NewHomeHandler(params HomeHandlerParams) *HomeHandler {
	return &HomeHandler{
		taskStore: params.TaskStore,
        sessionStore: params.SessionStore,
	}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    sessionID := r.Context().Value(middleware.SessionKey)

    if sessionID != nil {
        tasks , _:= h.taskStore.GetTasksFromSession(sessionID.(string))
        c := templates.Home(tasks)
        err := templates.HomeLayout(c, "ToDo App").Render(r.Context(), w)
        if err != nil { BadRequest(w, r); return}
        return
    }
    // If session doesnt exist yet, create one
    session, err := h.sessionStore.CreateSession()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	cookieValue := b64.StdEncoding.EncodeToString([]byte(session.ID))

    sessionCookieName := "session"

	expiration := time.Now().Add(365 * 24 * time.Hour)

	cookie := http.Cookie{
		Name:     sessionCookieName,
		Value:    cookieValue,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
    var tasks []store.Task
	http.SetCookie(w, &cookie)
    c := templates.Home(tasks)
    err = templates.HomeLayout(c, "ToDo App").Render(r.Context(), w)
    if err != nil { BadRequest(w, r); return}
}
