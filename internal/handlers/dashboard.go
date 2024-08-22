package handlers

import (
	"fmt"
	"net/http"

	"iteasy.wrappedAnsible/internal/model"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodGet); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	n := model.NewAnsibleProcessStatusDocument(r)
	result, err := n.Dashboard()

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to Query with DB: %v", err), http.StatusInternalServerError)
		return
	}

	_httpResponse(w, http.StatusOK, result)
}
