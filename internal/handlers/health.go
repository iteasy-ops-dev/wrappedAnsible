package handlers

import "net/http"

func health(w http.ResponseWriter, r *http.Request) {
	_httpResponse(w, http.StatusOK, nil)
}
