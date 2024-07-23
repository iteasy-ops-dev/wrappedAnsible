package handlers

import (
	"fmt"
	"net/http"
)

func AllowMethod(w http.ResponseWriter, r *http.Request, method string) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("invalid request method")
	}

	return nil
}
