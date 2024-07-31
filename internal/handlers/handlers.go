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
	mux.HandleFunc("/reset_password", ResetPassword)

	// with JWT
	mux.HandleFunc("/functions", Functions)
	mux.HandleFunc("/erp-parser", ErpParser)
	mux.HandleFunc("/run", ExcuteAnsible)
	mux.HandleFunc("/logs", Logs)
	mux.HandleFunc("/users", Users)
	mux.HandleFunc("/update_active", UpdateUserActive)
	mux.HandleFunc("/update_password", UpdateUserPassword)
}
