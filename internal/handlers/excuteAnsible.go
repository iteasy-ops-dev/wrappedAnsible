package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
	"iteasy.wrappedAnsible/internal/ansible"
)

// func ExcuteAnsible(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost {
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	ctx := r.Context()
// 	g, ctx := errgroup.WithContext(ctx)

// 	body, err := io.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Unable to read request body", http.StatusBadRequest)
// 		return
// 	}
// 	defer r.Body.Close()

// 	var payload []byte
// 	g.Go(func() error {
// 		var err error
// 		e := ansible.GetAnsibleFromFactory(ctx, body)
// 		payload, err = ansible.Excuter(e)
// 		return err
// 	})

// 	if err := g.Wait(); err != nil {
// 		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	_, err = w.Write(payload)
// 	if err != nil {
// 		http.Error(w, "Failed to write response", http.StatusInternalServerError)
// 		return
// 	}
// }

func ExcuteAnsible(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	g, ctx := errgroup.WithContext(ctx)

	p := ansible.GennerateHttpRequestType{
		Ctx: ctx,
		R:   r,
	}

	var payload []byte
	g.Go(func() error {
		e, err := ansible.GetAnsibleFromFactory(p)
		if err != nil {
			return fmt.Errorf("failed to get Ansible from factory: %w", err)
		}
		payload, err = ansible.Excuter(e)
		if err != nil {
			return fmt.Errorf("failed to execute Ansible: %w", err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
