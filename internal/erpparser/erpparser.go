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

// TODO: 현재 webhost manager에만 맞는 구조체임
// TODO: 확장이 필요함.
// TODO: 아래 코드 테스트 후 윗 주석 코드 날리면 됨
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

	data := url.Values{}
	data.Set("allow_type", config.GlobalConfig.Erp.Login.AllowType)
	data.Set("admin_id", config.GlobalConfig.Erp.Login.AdminId)
	data.Set("admin_passwd", config.GlobalConfig.Erp.Login.AdminPasswd)
	data.Set("login_btn", config.GlobalConfig.Erp.Login.LoginBtn)

	encodedData := data.Encode()

	loginReq, err := http.NewRequest("POST", config.GlobalConfig.Erp.Login.Url, strings.NewReader(encodedData))
	if err != nil {
		fmt.Println("HTTP 요청 생성 실패:", err)
		return
	}

	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	loginResp, err := client.Do(loginReq)
	if err != nil {
		fmt.Println("HTTP 요청 실패:", err)
		return
	}
	defer loginResp.Body.Close()

	var cookies []*http.Cookie
	cookies = append(cookies, jar.Cookies(&url.URL{Scheme: "https", Host: config.GlobalConfig.Erp.BaseURL})...)

	requestURL := e.Url

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		fmt.Println("HTTP 요청 생성 실패:", err)
		return
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("HTTP 요청 실패:", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("HTML 파싱 실패:", err)
		return
	}

	var combinedText string
	doc.Find("div.from-customer").Each(func(i int, s *goquery.Selection) {
		htmlContent, _ := s.Html()
		htmlContent = strings.ReplaceAll(htmlContent, "</p>", "</p>\n")
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
		if err != nil {
			fmt.Println("문서 로딩 오류:", err)
			return
		}
		text := doc.Text()
		combinedText += text + "\n"
	})
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
}

func (e *ErpParser) ToBytes() []byte {
	b, _ := json.Marshal(&e)
	return b
}

// TODO: 쓰리웨이 전용
// TODO: 범용적으로 활용 할 수 있는 방안을 찾아야 함.
// func extractInfo(text string) (Info, error) {
// 	info := Info{}

// 	patterns := map[string]*regexp.Regexp{
// 		"user_id":      regexp.MustCompile(`FTP 계정정보\s*ID\s*:\s*(\S+)`),
// 		"user_pass":    regexp.MustCompile(`PW\s*:\s*(\S+)`),
// 		"disk_quota":   regexp.MustCompile(`디스크할당량\s*:\s*(\d+)`),
// 		"cband_limit":  regexp.MustCompile(`1일 트래픽 할당량\s*:\s*(\d+)`),
// 		"vhost_domain": regexp.MustCompile(`연결도메인\s*:\s*(\S+)`),
// 		"db_user":      regexp.MustCompile(`DB 계정정보\s*ID\s*:\s*(\S+)`),
// 		"db_name":      regexp.MustCompile(`DBNAME\s*:\s*(\S+)`),
// 		"db_password":  regexp.MustCompile(`PW\s*:\s*(\S+)`),
// 	}

// 	for key, re := range patterns {
// 		match := re.FindStringSubmatch(text)
// 		if len(match) > 1 {
// 			switch key {
// 			case "user_id":
// 				info.UserID = match[1]
// 				fmt.Println("UserID:", info.UserID)
// 			case "user_pass":
// 				info.UserPass = match[1]
// 				fmt.Println("UserPass:", info.UserPass)
// 			case "disk_quota":
// 				info.DiskQuota = match[1]
// 				fmt.Println("DiskQuota:", info.DiskQuota)
// 			case "cband_limit":
// 				info.CbandLimit = match[1]
// 				fmt.Println("CbandLimit:", info.CbandLimit)
// 			case "vhost_domain":
// 				info.VhostDomain = match[1]
// 				fmt.Println("VhostDomain:", info.VhostDomain)
// 			case "db_user":
// 				info.DBUser = match[1]
// 				fmt.Println("DBUser:", info.DBUser)
// 			case "db_name":
// 				info.DBName = match[1]
// 				fmt.Println("DBName:", info.DBName)
// 			case "db_password":
// 				info.DBPassword = match[1]
// 				fmt.Println("DBPassword:", info.DBPassword)
// 			}
// 		}
// 	}

// 	return info, nil
// }

// func extractInfo(text string) (Info, error) {
// 	info := Info{}

// 	// 정규 표현식 패턴 정의
// 	patterns := map[string]*regexp.Regexp{
// 		"user_id":      regexp.MustCompile(`FTP 계정정보\s*ID\s*:\s*(\S+)`),
// 		"user_pass":    regexp.MustCompile(`PW\s*:\s*(\S+)`),
// 		"disk_quota":   regexp.MustCompile(`디스크할당량\s*:\s*(\d+)`),
// 		"cband_limit":  regexp.MustCompile(`1일 트래픽 할당량\s*:\s*(\d+)`),
// 		"vhost_domain": regexp.MustCompile(`연결도메인\s*:\s*(\S+)`),
// 		"db_user":      regexp.MustCompile(`DB 계정정보\s*ID\s*:\s*(\S+)`),
// 		"db_name":      regexp.MustCompile(`DBNAME\s*:\s*(\S+)`),
// 		// "db_password" 패턴 제거: 두 개의 PW 값을 찾기 위해 특수 처리를 합니다.
// 	}

// 	// 먼저 일반 패턴들을 처리합니다.
// 	for key, re := range patterns {
// 		match := re.FindStringSubmatch(text)
// 		if len(match) > 1 {
// 			switch key {
// 			case "user_id":
// 				info.UserID = match[1]
// 			case "user_pass":
// 				info.UserPass = match[1]
// 			case "disk_quota":
// 				info.DiskQuota = match[1]
// 			case "cband_limit":
// 				info.CbandLimit = match[1]
// 			case "vhost_domain":
// 				info.VhostDomain = match[1]
// 			case "db_user":
// 				info.DBUser = match[1]
// 			case "db_name":
// 				info.DBName = match[1]
// 			}
// 		}
// 	}

// 	// "PW" 패턴은 FindAllStringSubmatch로 처리합니다.
// 	passMatches := regexp.MustCompile(`PW\s*:\s*(\S+)`).FindAllStringSubmatch(text, -1)
// 	if len(passMatches) >= 2 {
// 		info.UserPass = passMatches[0][1]   // 첫 번째 PW 값
// 		info.DBPassword = passMatches[1][1] // 두 번째 PW 값
// 		// fmt.Println("UserPass:", info.UserPass)
// 		// fmt.Println("DBPassword:", info.DBPassword)
// 	}

// 	return info, nil
// }

// 앞뒤 공백 제거
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
	}

	// 먼저 일반 패턴들을 처리합니다.
	for key, re := range patterns {
		match := re.FindStringSubmatch(text)
		if len(match) > 1 {
			trimmedValue := strings.TrimSpace(match[1])
			switch key {
			case "user_id":
				info.UserID = trimmedValue
			case "user_pass":
				info.UserPass = trimmedValue
			case "disk_quota":
				info.DiskQuota = trimmedValue
			case "cband_limit":
				info.CbandLimit = trimmedValue
			case "vhost_domain":
				info.VhostDomain = trimmedValue
			case "db_user":
				info.DBUser = trimmedValue
			case "db_name":
				info.DBName = trimmedValue
			}
		}
	}

	// "PW" 패턴은 FindAllStringSubmatch로 처리합니다.
	passMatches := regexp.MustCompile(`PW\s*:\s*(\S+)`).FindAllStringSubmatch(text, -1)
	if len(passMatches) >= 2 {
		info.UserPass = strings.TrimSpace(passMatches[0][1])   // 첫 번째 PW 값
		info.DBPassword = strings.TrimSpace(passMatches[1][1]) // 두 번째 PW 값
	}

	return info, nil
}
