package erpparser

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	config "iteasy.wrappedAnsible/configs"

	"github.com/PuerkitoBio/goquery"
)

type Info struct {
	UserID      string
	UserPass    string
	DiskQuota   string
	CbandLimit  string
	VhostDomain string
	DBUser      string
	DBName      string
	DBPassword  string
}

type ErpParser struct {
	Url  string
	Info Info
}

func New(url string) *ErpParser {
	return &ErpParser{
		Url: url,
	}
}

func (e *ErpParser) Parsing() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("쿠키 저장소 생성 실패:", err)
		return
	}

	// 데이터 정의
	// 로그인 요청에 필요한 데이터를 URL 인코딩된 형식으로 설정합니다.
	data := url.Values{}
	data.Set("allow_type", config.GlobalConfig.Erp.Login.AllowType)
	data.Set("admin_id", config.GlobalConfig.Erp.Login.AdminId)
	data.Set("admin_passwd", config.GlobalConfig.Erp.Login.AdminPasswd)
	data.Set("login_btn", config.GlobalConfig.Erp.Login.LoginBtn)

	// URL 인코딩된 데이터 준비
	// `url.Values`로 생성된 데이터를 URL 인코딩된 문자열로 변환합니다.
	encodedData := data.Encode()

	// 요청 생성
	// 로그인 요청을 POST 방식으로 생성합니다.
	loginReq, err := http.NewRequest("POST", config.GlobalConfig.Erp.Login.Url, strings.NewReader(encodedData))
	// loginReq, err := http.NewRequest("POST", "https://admin.ksidc.net/?act=login&ref=%2F", strings.NewReader(encodedData))
	if err != nil {
		fmt.Println("HTTP 요청 생성 실패:", err)
		return
	}

	// Content-Type 헤더 설정
	// 요청 헤더를 설정하여 서버가 이해할 수 있는 형식으로 데이터를 전송합니다.
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 클라이언트 생성
	// 쿠키 저장소를 클라이언트에 설정하고 리다이렉트를 처리할 방법을 설정합니다.
	client := &http.Client{
		Jar: jar, // 쿠키 저장소 설정
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 리다이렉트가 발생할 때마다 쿠키를 추적하기 위해, 리다이렉트를 허용합니다.
			return nil
		},
	}

	// 로그인 요청 보내기
	// 설정한 클라이언트를 사용하여 로그인 요청을 서버로 전송합니다.
	loginResp, err := client.Do(loginReq)
	if err != nil {
		fmt.Println("HTTP 요청 실패:", err)
		return
	}
	defer loginResp.Body.Close() // 응답 본문을 다 읽은 후 연결을 닫습니다.

	// 로그인 응답 상태 코드 출력
	// 로그인 요청에 대한 응답 상태 코드를 출력하여 성공 여부를 확인합니다.
	fmt.Println("로그인 응답 상태 코드:", loginResp.StatusCode)

	// 최종 응답에서 쿠키 읽기
	// 응답에서 쿠키를 읽어와 출력합니다.
	if len(loginResp.Cookies()) == 0 {
		fmt.Println("응답에서 쿠키를 찾을 수 없습니다.")
	} else {
		for _, cookie := range loginResp.Cookies() {
			fmt.Printf("쿠키 이름: %s, 쿠키 값: %s\n", cookie.Name, cookie.Value)
		}
	}

	// 쿠키 저장소의 쿠키 출력
	// 쿠키 저장소에 저장된 모든 쿠키를 출력하여, 리다이렉트 후 설정된 쿠키를 확인합니다.
	var cookies []*http.Cookie
	fmt.Println("저장된 쿠키:")
	for _, cookie := range jar.Cookies(&url.URL{Scheme: "https", Host: config.GlobalConfig.Erp.BaseURL}) {
		fmt.Printf("쿠키 이름: %s, 쿠키 값: %s\n", cookie.Name, cookie.Value)
		cookies = append(cookies, cookie)
	}

	// 2단계: 얻은 쿠키를 사용하여 다른 페이지 요청 보내기
	// requestURL := "https://admin.ksidc.net/service/request_info/?mem_idx=40396&idx=224841"
	// requestURL := "https://admin.ksidc.net/service/request_info/?mem_idx=40396&idx=225737"
	requestURL := e.Url

	// 새로운 요청 생성
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println("HTTP 요청 생성 실패:", err)
		return
	}

	// 쿠키 설정
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP 요청 실패:", err)
		return
	}
	defer resp.Body.Close()

	// 응답 상태 코드 출력
	fmt.Println("응답 상태 코드:", resp.StatusCode)

	// HTML 파싱
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("HTML 파싱 실패:", err)
		return
	}

	// 원하는 div에서 텍스트 추출
	var combinedText string
	doc.Find("div.from-customer").Each(func(i int, s *goquery.Selection) {
		htmlContent, _ := s.Html()
		// <p> 태그를 줄바꿈으로 대체
		htmlContent = strings.ReplaceAll(htmlContent, "</p>", "</p>\n")
		// HTML을 텍스트로 변환 (태그 제거)
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			fmt.Println("문서 로딩 오류:", err)
			return
		}
		// 각 <p> 태그 사이에 줄바꿈 추가
		text := doc.Text()
		// 줄바꿈 추가
		combinedText += text + "\n"
	})
	// Trim any extra spaces at the end
	combinedText = strings.TrimSpace(combinedText)

	info, err := extractInfo(combinedText)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	e.Info.UserID = info.UserID
	e.Info.UserPass = info.UserPass
	e.Info.DiskQuota = info.DiskQuota
	e.Info.CbandLimit = info.CbandLimit
	e.Info.VhostDomain = info.VhostDomain
	e.Info.DBUser = info.DBUser
	e.Info.DBName = info.DBName
	e.Info.DBPassword = info.DBPassword

	// fmt.Println("Extracted Information:")
	// fmt.Printf("UserID: %s\n", info.UserID)
	// fmt.Printf("UserPass: %s\n", info.UserPass)
	// fmt.Printf("DiskQuota: %s\n", info.DiskQuota)
	// fmt.Printf("CbandLimit: %s\n", info.CbandLimit)
	// fmt.Printf("VhostDomain: %s\n", info.VhostDomain)
	// fmt.Printf("DBUser: %s\n", info.DBUser)
	// fmt.Printf("DBName: %s\n", info.DBName)
	// fmt.Printf("DBPassword: %s\n", info.DBPassword)
}

func (e *ErpParser) ToBytes() []byte {
	b, _ := json.Marshal(&e)
	return b
}

// TODO: 작업의뢰를 세분화해서 가져 올 수 있게 하는 로직이 필요
func extractInfo(text string) (Info, error) {
	info := Info{}

	// 정규 표현식 패턴 정의
	patterns := map[string]*regexp.Regexp{
		"user_id":      regexp.MustCompile(`FTP 계정정보\s*ID\s*:\s*(\S+)`),
		"user_pass":    regexp.MustCompile(`PW\s*:\s*(\S+)`),
		"disk_quota":   regexp.MustCompile(`디스크할당량\s*:\s*(\d+)`),
		"cband_limit":  regexp.MustCompile(`1일 트래픽 할당량\s*:\s*(\d+)`),
		"vhost_domain": regexp.MustCompile(`연결도메인\s*:\s*(\S+)`),
		"db_user":      regexp.MustCompile(`DB 계정정보\s*ID\s*:\s*(\S+)`),
		"db_name":      regexp.MustCompile(`DBNAME\s*:\s*(\S+)`),
		"db_password":  regexp.MustCompile(`PW\s*:\s*(\S+)`),
	}
	// 정규 표현식으로 정보 추출
	for key, re := range patterns {
		match := re.FindStringSubmatch(text)
		if len(match) > 1 {
			switch key {
			case "user_id":
				info.UserID = match[1]
			case "user_pass":
				info.UserPass = match[1]
			case "disk_quota":
				info.DiskQuota = match[1]
			case "cband_limit":
				info.CbandLimit = match[1]
			case "vhost_domain":
				info.VhostDomain = match[1]
			case "db_user":
				info.DBUser = match[1]
			case "db_name":
				info.DBName = match[1]
			case "db_password":
				info.DBPassword = match[1]
			}
		}
	}

	return info, nil
}
