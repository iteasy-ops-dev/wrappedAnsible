package handlers

import (
	"net/http"
)

// RegisterHandlers registers all the route handlers
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/health", Health)
	mux.HandleFunc("/signup", SignUp)
	mux.HandleFunc("/verify", VerifyEmail)
	mux.HandleFunc("/login", Login)
	mux.HandleFunc("/logout", Logout)

	// with JWT
	mux.HandleFunc("/functions", Functions)
	mux.HandleFunc("/run", ExcuteAnsible)
	mux.HandleFunc("/get", Get)
	mux.HandleFunc("/erp-parser", ErpParser)
}
