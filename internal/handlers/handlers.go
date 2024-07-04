package handlers

import (
	"net/http"
)

// RegisterHandlers registers all the route handlers
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/list_func", ListFunc)
	mux.HandleFunc("/excuteAnsible", ExcuteAnsible)
	mux.HandleFunc("/excuteAnsibleWithFiles", ExcuteAnsibleWithFiles)
}
