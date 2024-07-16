package handlers

import (
	"net/http"
)

// RegisterHandlers registers all the route handlers
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/signup", SignUp)
	mux.HandleFunc("/login", Login)

	// with JWT
	mux.HandleFunc("/list_func", ListFunc)
	mux.HandleFunc("/run", ExcuteAnsible)
	mux.HandleFunc("/get", Get)
}
