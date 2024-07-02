package router

import (
	"net/http"

	"iteasy.wrappedAnsible/internal/handlers"
)

// NewRouter creates a new HTTP router and registers all the handlers
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Registering all handlers
	handlers.RegisterHandlers(mux)

	return mux
}
