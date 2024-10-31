package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/sync/errgroup"

	"iteasy.wrappedAnsible/internal/ansible"
	"iteasy.wrappedAnsible/internal/model"
)

func excuteAnsible(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// TODO: 확장자 검사 추가. 화이트리스트로 검증
		err := r.ParseMultipartForm(10 << 20) // 10MB max memory
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Failed to parse form data", http.StatusBadRequest)
			return
		}
	}

	g, ctx := errgroup.WithContext(r.Context())

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
		m := model.NewAnsibleProcessStatusDocument(payload)
		if err := m.Put(); err != nil {
			return err
		}
		return err
	})

	if err := g.Wait(); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	_httpResponse(w, http.StatusOK, payload)
}
