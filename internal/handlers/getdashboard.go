package handlers

import (
	// "goth/internal/middleware"
	// "goth/internal/store"
	// "goth/internal/templates"
	// "net/http"
)

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

// func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//
// 	user, ok := r.Context().Value(middleware.UserKey).(*store.User)
// 	if !ok {
// 		c := templates.GuestIndex()
// 		err := templates.DashboardLayout(c, "My website").Render(r.Context(), w)
//
// 		if err != nil {
// 			http.Error(w, "Error rendering template", http.StatusInternalServerError)
// 			return
// 		}
//
// 		return
// 	}
//
// 	c := templates.Index(user.Email)
// 	err := templates.DashboardLayout(c, "My website").Render(r.Context(), w)
//
// 	if err != nil {
// 		http.Error(w, "Error rendering template", http.StatusInternalServerError)
// 		return
// 	}
// }

