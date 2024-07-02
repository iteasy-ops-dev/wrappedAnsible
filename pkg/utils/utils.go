package utils

import (
	"log"
	"os"
	"regexp"
	"strings"
)

func AddEmptySpace(data string) string {
	re := regexp.MustCompile(`\n`)
	return re.ReplaceAllString(data, "\n          ")
}

// 파일 사이즈 가져오기
func GetFileSize(file *os.File) int64 {
	stat, err := os.Stat(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	return stat.Size()
}

// 파일 존재 확인
func ExistFile(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

// 파일 내용 전체 읽기
func ReadFile(file *os.File) []byte {
	f, err := os.Open(file.Name()) // For read access.
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, GetFileSize(f))
	count, err := file.Read(data)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("read %d bytes: %q\n", count, data[:count])
	return data[:count]
}

// 임시파일 생성하기
func GenerateTempFile(payload []byte, namePattern string) (*os.File, error) {
	// 파일 생성 부분
	tmp, err := os.CreateTemp("", namePattern)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmp.Write(payload); err != nil {
		log.Fatal(err)
	}
	if err := tmp.Close(); err != nil {
		log.Fatal(err)
	}

	return tmp, nil
}

// 폴더 내 파일 리스트 받아오기
func GetFileList(path string) []string {
	var list []string
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		list = append(list, file.Name())
	}

	return list
}

// 파일 리스트를 Map타입으로
func GetFileListForMap(path string) map[string]string {
	mlist := make(map[string]string)

	list := GetFileList(path)
	for _, filename := range list {
		mlist[TruncationExtension(filename)] = filename
	}
	return mlist
}

// 확장자 자르기
func TruncationExtension(filename string) string {
	var s string
	if ExistFile(filename) {
		return ""
	}
	dotIndex := strings.LastIndex(filename, ".")
	if dotIndex != -1 {
		s = filename[:dotIndex]
	} else {
		s = filename
	}
	return s
}

// 확장자 체크
func CheckExtension(filename, pattern string) bool {
	// pattern이 확장자인지 확인하는 절차 필요
	return strings.Contains(filename, pattern)
}
