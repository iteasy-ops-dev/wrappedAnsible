package handlers

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/sync/errgroup"
	"iteasy.wrappedAnsible/internal/ansible"
)

func GetParam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	g, ctx := errgroup.WithContext(ctx)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var payload []byte
	g.Go(func() error {
		var err error
		e := ansible.GetAnsibleFromFactory(ctx, body)
		payload, err = ansible.Excuter(e)
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(payload)
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

// func GetParam(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Unable to read request body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	e := ansible.GetAnsibleFromFactory(body)
// 	payload := ansible.Excuter(e)

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	_, err = w.Write(payload)
// 	if err != nil {
// 		http.Error(w, "Failed to write response", http.StatusInternalServerError)
// 		return
// 	}
// }
