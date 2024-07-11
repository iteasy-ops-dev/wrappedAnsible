package ansible

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
)

func generateInitAnsible(v GennerateInitType) (iAnsible, error) {
	e := extendAnsible{
		Ctx: v.Ctx,
	}

	if err := json.Unmarshal(v.JsonData, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func generateHttpAnsible(v GennerateHttpRequestType) (iAnsible, error) {
	contentType := v.R.Header.Get("Content-Type")
	switch {
	case strings.Contains(contentType, "application/json"):
		return parseJsonRequest(v)
	case strings.Contains(contentType, "multipart/form-data"):
		return parseMultipartRequest(v)
	default:
		return nil, errors.New("구성 할 수 없는 Content-type")
	}
}

func parseJsonRequest(v GennerateHttpRequestType) (iAnsible, error) {
	body, err := io.ReadAll(v.R.Body)
	if err != nil {
		return nil, errors.New("unable to read request body")
	}
	defer v.R.Body.Close()

	e := extendAnsible{
		Ctx: v.Ctx,
	}

	if err := json.Unmarshal(body, &e); err != nil {
		return nil, err
	}

	return &e, nil
}

func parseMultipartRequest(v GennerateHttpRequestType) (iAnsible, error) {
	e := extendAnsible{
		Ctx:     v.Ctx,
		Options: make(map[string]string),
	}

	// Form values
	for key, values := range v.R.Form {
		if len(values) > 0 {
			if err := _setField(&e, key, values[0]); err != nil {
				return nil, err
			}
		}
	}

	// Handling file uploads
	files := v.R.MultipartForm.File["files"]
	var filePaths []string
	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		// Create a temporary file
		tempFile, err := os.CreateTemp("", fileHeader.Filename)
		if err != nil {
			return nil, err
		}
		defer tempFile.Close()

		// Read the file content and write it to the temporary file
		if _, err := io.Copy(tempFile, file); err != nil {
			return nil, err
		}

		filePaths = append(filePaths, tempFile.Name())
	}

	for _, fp := range filePaths {
		e.Options[_generateOptionFileKey(fp, e.Type)] = fp
	}

	return &e, nil
}

// ===== 해당 파일 안에서만 사용하는 함수 =====
func _setField(e *extendAnsible, key, value string) error {
	switch key {
	case "type":
		e.Type = value
	case "name":
		e.Name = value
	case "account":
		e.Account = value
	case "password":
		e.Password = value
	case "description":
		e.Description = value
	case "ips":
		if err := json.Unmarshal([]byte(value), &e.IPs); err != nil {
			return err
		}
	default:
		// Handle additional fields if needed
		e.Options[key] = value
	}
	return nil
}

// Options.src_키_file 생성기
// ansible roles 와의 파일 부분에 대한 옵션을 맞춰주는 함수
// TODO: 에러 제어 필요
func _generateOptionFileKey(filepath, extendAnsible_Type string) string {
	switch extendAnsible_Type {
	case "change_ssl":
		switch {
		case strings.Contains(filepath, "key"):
			return "src_key_file"
		case strings.Contains(filepath, "crt"):
			return "src_cert_file"
		case strings.Contains(filepath, "chain"):
			return "src_chain_file"
		default:
			return ""
		}

	default:
		return ""
	}
}
