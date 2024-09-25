package handlers

import (
	"net/http"
)

// RegisterHandlers registers all the route handlers
func RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/health", health)
	mux.HandleFunc("/signup", signUp)
	mux.HandleFunc("/verify", verifyEmail)
	mux.HandleFunc("/login", login)
	mux.HandleFunc("/logout", logout)
	mux.HandleFunc("/reset_password", resetPassword)

	mux.HandleFunc("/ws", handleWebSocket)

	// with JWT
	mux.HandleFunc("/extend_extension", extendJWT)
	mux.HandleFunc("/functions", functions)
	mux.HandleFunc("/dashboad", dashboard)
	mux.HandleFunc("/run", excuteAnsible)
	mux.HandleFunc("/logs", logs)
	mux.HandleFunc("/users", users)
	mux.HandleFunc("/update_active", updateUserActive)
	mux.HandleFunc("/update_password", updateUserPassword)

	mux.HandleFunc("/erp-parser", erpParser)
	mux.HandleFunc("/update-work-history", updateWorkHistory)
	mux.HandleFunc("/workhistory", workHistory)
	mux.HandleFunc("/get-work-history", getWorkHistory)
	mux.HandleFunc("/insight", insight)
}
