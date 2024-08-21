package handlers

import (
	"fmt"
	"net/http"

	"iteasy.wrappedAnsible/internal/model"
)

func logs(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodGet); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}
	n := model.NewAnsibleProcessStatusDocument(r)
	result, err := n.Get()

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to Query with DB: %v", err), http.StatusInternalServerError)
		return
	}

	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)
	// if err := json.NewEncoder(w).Encode(result); err != nil {
	// 	http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	// 	return
	// }

	_httpResponse(w, http.StatusOK, result)
}
