package utils

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestAddEmptySpace(t *testing.T) {
	data := "Hello\nWorld"
	expected := "Hello\n          World"
	result := AddEmptySpace(data)
	if result != expected {
		t.Errorf("AddEmptySpace() = %q; want %q", result, expected)
	}
}

func TestGetFileSize(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	_, err = tmpfile.Write([]byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	expectedSize := int64(4)
	result := GetFileSize(tmpfile)
	if result != expectedSize {
		t.Errorf("GetFileSize() = %d; want %d", result, expectedSize)
	}
}

func TestExistFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	exists := ExistFile(tmpfile.Name())
	if !exists {
		t.Errorf("ExistFile() = false; want true")
	}

	exists = ExistFile("nonexistent_file")
	if exists {
		t.Errorf("ExistFile() = true; want false")
	}
}

func TestReadFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	content := []byte("test content")
	_, err = tmpfile.Write(content)
	if err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	file, err := os.Open(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	result := ReadFile(file)
	if !bytes.Equal(result, content) {
		t.Errorf("ReadFile() = %q; want %q", result, content)
	}
}

func TestGenerateTempFile(t *testing.T) {
	payload := []byte("temporary content")
	tmpfile, err := GenerateTempFile(payload, "testfile_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up

	result := ReadFile(tmpfile)
	if !bytes.Equal(result, payload) {
		t.Errorf("GenerateTempFile() = %q; want %q", result, payload)
	}
}

func TestGetFileList(t *testing.T) {
	dir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // Clean up

	files := []string{"file1.txt", "file2.txt"}
	for _, file := range files {
		_, err := os.Create(dir + "/" + file)
		if err != nil {
			t.Fatal(err)
		}
	}

	result := GetFileList(dir)
	if len(result) != len(files) {
		t.Errorf("GetFileList() length = %d; want %d", len(result), len(files))
	}
}

func TestTruncationExtension(t *testing.T) {
	filename := "file.txt"
	expected := "file"
	result := TruncationExtension(filename)
	if result != expected {
		t.Errorf("TruncationExtension() = %q; want %q", result, expected)
	}
}

func TestCheckExtension(t *testing.T) {
	filename := "file.txt"
	pattern := ".txt"
	if !CheckExtension(filename, pattern) {
		t.Errorf("CheckExtension() = false; want true")
	}

	pattern = ".jpg"
	if CheckExtension(filename, pattern) {
		t.Errorf("CheckExtension() = true; want false")
	}
}

func TestGenerateToken(t *testing.T) {
	token := GenerateToken()
	if len(token) != 32 {
		t.Errorf("GenerateToken() length = %d; want 32", len(token))
	}
}

func TestGenerateTempPassword(t *testing.T) {
	password, err := GenerateTempPassword()
	if err != nil {
		t.Fatal(err)
	}

	if len(password) < 12 {
		t.Errorf("GenerateTempPassword() length = %d; want >= 12", len(password))
	}

	_, err = base64.URLEncoding.DecodeString(password)
	if err != nil {
		t.Errorf("GenerateTempPassword() decoding failed: %v", err)
	}
}

func TestHashingPassword(t *testing.T) {
	password := "password123"
	hash, err := HashingPassword(password)
	if err != nil {
		t.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		t.Errorf("HashingPassword() hash comparison failed: %v", err)
	}
}

func TestParseRequestBody(t *testing.T) {
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	body := `{"name": "test", "value": 123}`
	req := &http.Request{
		Body: io.NopCloser(bytes.NewBufferString(body)),
	}

	result, err := ParseRequestBody[TestData](req)
	if err != nil {
		t.Fatal(err)
	}

	if result.Name != "test" || result.Value != 123 {
		t.Errorf("ParseRequestBody() = %+v; want %+v", result, TestData{Name: "test", Value: 123})
	}
}
