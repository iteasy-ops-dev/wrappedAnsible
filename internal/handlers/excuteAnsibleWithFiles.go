package handlers

import (
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
	"iteasy.wrappedAnsible/internal/ansible"
)

// TODO: FormData로 들어온 오브젝트를
// extendAnsible 타입에 맞게 변환해야한다.
// 파일은 임시저장하고
// 해당 경로를 options의 하위에 저장한다.
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

	fmt.Println(r.Header.Get("Content-Type"))

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

// func ExcuteAnsibleWithFiles(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodPost { // MethodPost
// 		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	err := r.ParseMultipartForm(10 << 20) // 10MB max memory
// 	if err != nil {
// 		http.Error(w, "Unable to parse form", http.StatusBadRequest)
// 		return
// 	}

// 	if err := r.ParseForm(); err != nil {
// 		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
// 		return
// 	}

// 	name := r.FormValue("name")
// 	account := r.FormValue("account")
// 	password := r.FormValue("password")
// 	domain := r.FormValue("domain")

// 	ipJSON := r.FormValue("ips")
// 	var ips []string
// 	if err := json.Unmarshal([]byte(ipJSON), &ips); err != nil {
// 		http.Error(w, "Failed to parse IP JSON data", http.StatusBadRequest)
// 		return
// 	}

// 	// Variables to store file contents
// 	var certContent, chainContent, keyContent string

// 	// Loop through the uploaded files
// 	for _, fileHeader := range r.MultipartForm.File["files"] {
// 		file, err := fileHeader.Open()
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error opening file %s", fileHeader.Filename), http.StatusInternalServerError)
// 			return
// 		}
// 		defer file.Close()

// 		content, err := io.ReadAll(file)
// 		if err != nil {
// 			http.Error(w, fmt.Sprintf("Error reading file %s", fileHeader.Filename), http.StatusInternalServerError)
// 			return
// 		}

// 		// Assign content based on the file name
// 		if strings.Contains(fileHeader.Filename, "crt") {
// 			certContent = string(content)
// 		} else if strings.Contains(fileHeader.Filename, "chain") {
// 			chainContent = string(content)
// 		} else if strings.Contains(fileHeader.Filename, "key") {
// 			keyContent = string(content)
// 		}
// 	}

// 	fmt.Println(certContent)
// 	fmt.Println(chainContent)
// 	fmt.Println(keyContent)
// 	fmt.Println(name)
// 	fmt.Println(account)
// 	fmt.Println(password)
// 	fmt.Println(domain)
// 	fmt.Println(ipJSON)
// 	// Send a response
// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte("Upload successful"))
// }
