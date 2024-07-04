package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

	name := r.FormValue("name")
	account := r.FormValue("account")
	password := r.FormValue("password")
	domain := r.FormValue("domain")

	ipJSON := r.FormValue("ips")
	var ips []string
	if err := json.Unmarshal([]byte(ipJSON), &ips); err != nil {
		http.Error(w, "Failed to parse IP JSON data", http.StatusBadRequest)
		return
	}

	// Variables to store file contents
	var certContent, chainContent, keyContent string

	// Loop through the uploaded files
	for _, fileHeader := range r.MultipartForm.File["files"] {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, fmt.Sprintf("Error opening file %s", fileHeader.Filename), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading file %s", fileHeader.Filename), http.StatusInternalServerError)
			return
		}

		// Assign content based on the file name
		if strings.Contains(fileHeader.Filename, "crt") {
			certContent = string(content)
		} else if strings.Contains(fileHeader.Filename, "chain") {
			chainContent = string(content)
		} else if strings.Contains(fileHeader.Filename, "key") {
			keyContent = string(content)
		}
	}

	fmt.Println(certContent)
	fmt.Println(chainContent)
	fmt.Println(keyContent)
	fmt.Println(name)
	fmt.Println(account)
	fmt.Println(password)
	fmt.Println(domain)
	fmt.Println(ipJSON)
	// Send a response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Upload successful"))
}
