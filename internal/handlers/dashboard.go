package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"iteasy.wrappedAnsible/internal/model"
	"iteasy.wrappedAnsible/pkg/utils"
)

func dashboard(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	var req struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	n := model.NewAnsibleProcessStatusDocument(r)
	result, err := n.Dashboard(utils.StartOfDay(req.Start), utils.EndOfDay(req.End))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to Query with DB: %v", err), http.StatusInternalServerError)
		return
	}

	_httpResponse(w, http.StatusOK, result)
}
