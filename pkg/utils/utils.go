package utils

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/rand"
	config "iteasy.wrappedAnsible/configs"
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
	// log.Printf("read %d bytes: %q\n", count, data[:count])
	return data[:count]
}

// 임시파일 생성하기
func GenerateTempFile(payload []byte, namePattern string) (*os.File, error) {
	// 파일 생성 부분
	tmp, err := os.CreateTemp("", namePattern)
	if err != nil {
		log.Fatal(err)
	}
	defer tmp.Close()

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

func RemoveFile(path string) error {
	return os.Remove(path)
}

// 메일 보내기
func SendEmail(to, subject, mailBody string) error {
	// 메일 헤더와 본문 작성
	// msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", config.GlobalConfig.Smtp.From, to, subject, mailBody)
	msg := "MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		fmt.Sprintf("From: %s\r\n", config.GlobalConfig.Smtp.From) +
		fmt.Sprintf("To: %s\r\n", to) +
		fmt.Sprintf("Subject: %s\r\n", subject) +
		"\r\n" + mailBody

	// 메일 전송
	err := smtp.SendMail(
		config.GlobalConfig.Smtp.Host+":"+config.GlobalConfig.Smtp.Port,
		nil, // Auth
		config.GlobalConfig.Smtp.From,
		[]string{to},
		[]byte(msg),
	)

	if err != nil {
		log.Printf("메일 전송 실패: %v\n", err)
		return err
	}
	return nil
}

func GenerateToken() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func GenerateTempPassword() (string, error) {
	length := 12
	randBytes := make([]byte, length)
	_, err := rand.Read(randBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randBytes), nil
}

func HashingPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func ParseRequestBody[T any](r *http.Request) (T, error) {
	var body T
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return body, fmt.Errorf("failed to read request body: %v", err)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &body); err != nil {
		return body, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return body, nil
}

func ParseJSONBody(body io.Reader, v interface{}) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return nil
}

func VerifySSL(file string) (string, error) {
	// 인증서 파일 읽기
	certPEM, err := os.ReadFile(file)
	if err != nil {
		// log.Printf("인증서 파일을 읽는 중 오류 발생: %v\n", err)
		return "", err
	}

	// PEM 블록 디코드
	// block, rest := pem.Decode(certPEM)
	block, _ := pem.Decode(certPEM)
	if block == nil {
		// fmt.Println("유효한 PEM 데이터가 없습니다.")
		return "", err
	}
	// if len(rest) > 0 {
	// 	fmt.Println("경고: PEM 블록 뒤에 추가 데이터가 발견되었습니다.")
	// }

	// 인증서 파싱
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		// log.Printf("=================== 키 파일: %s\n", file)
		// log.Printf("인증서 파싱 중 오류 발생: %v\n", err)
		return "key", nil
	}

	// 인증서 정보 출력
	// fmt.Println("주체:", cert.Subject)
	// fmt.Println("발급자:", cert.Issuer)
	// fmt.Println("유효 시작일:", cert.NotBefore)
	// fmt.Println("유효 종료일:", cert.NotAfter)
	// fmt.Println("일련 번호:", cert.SerialNumber)
	// fmt.Println("DNS 이름들:", cert.DNSNames)

	// CA 인증서 여부 확인
	if cert.IsCA {
		// fmt.Println("이 인증서는 CA 인증서입니다.")
		// log.Printf("=================== ca 파일: %s\n", file)
		return "ca", nil
	} else {
		// fmt.Println("이 인증서는 CA 인증서가 아닙니다.")
		// log.Printf("=================== crt 파일: %s\n", file)
		return "crt", nil
	}
}
