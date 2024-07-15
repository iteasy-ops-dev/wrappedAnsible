package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
	"iteasy.wrappedAnsible/internal/ansible"
	"iteasy.wrappedAnsible/internal/model"
)

func ExcuteAnsibleWithFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost { // MethodPost
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10MB max memory
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	// fmt.Println("Content-Type: ")
	// fmt.Println(r.Header.Get("Content-Type"))

	ctx := r.Context()
	g, ctx := errgroup.WithContext(ctx)

	p := ansible.GennerateHttpRequestType{
		Ctx: ctx,
		R:   r,
	}

	var payload *ansible.AnsibleProcessStatus
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

	m := model.NewAnsibleProcessStatus(payload)
	m.Put()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(payload.ToBytes()); err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}
