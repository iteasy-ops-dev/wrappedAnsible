package erpparser

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/pkg/utils"
)

// WorkHistory 구조체 정의
type WorkHistory struct {
	Index                  string
	ProcessingStatus       string
	RegistrationDate       string
	DesiredWorkTime        string
	ExpectedCompletionTime string
	Worker                 string
	RequestItemName        string
	RequestContent         string
	RequestLink            string
	Subcategory            string
	CompanyName            string
	IP                     string
	Brand                  string
	RegistrationType       string
}

// 날짜를 처리하는 함수
func WorkHistoryParsing() ([]WorkHistory, error) {
	// CookieJar 생성
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("쿠키 저장소 생성 실패: %v", err)
	}

	// 로그인 요청 데이터 설정
	data := url.Values{}
	data.Set("allow_type", config.GlobalConfig.Erp.Login.AllowType)
	data.Set("admin_id", config.GlobalConfig.Erp.Login.AdminId)
	data.Set("admin_passwd", config.GlobalConfig.Erp.Login.AdminPasswd)
	data.Set("login_btn", config.GlobalConfig.Erp.Login.LoginBtn)

	// 로그인 요청 생성
	loginReq, err := http.NewRequest("POST", config.GlobalConfig.Erp.Login.Url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 생성 실패: %v", err)
	}
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// HTTP 클라이언트 생성 (쿠키 포함)
	client := &http.Client{
		Jar: jar,
	}

	// 로그인 요청 보내기
	loginResp, err := client.Do(loginReq)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 실패: %v", err)
	}
	defer loginResp.Body.Close()

	// 로그인 성공 여부 확인
	if loginResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("로그인 실패: 상태 코드 %d", loginResp.StatusCode)
	}

	payload := url.Values{
		"0[name]":    {"hopeCal1"},
		"0[value]":   {"0000-00-00"},
		"1[name]":    {"hopeCal2"},
		"1[value]":   {"0000-00-00"},
		"2[name]":    {"Cal1"},
		"2[value]":   {"0000-00-00"},
		"3[name]":    {"Cal2"},
		"3[value]":   {"0000-00-00"},
		"4[name]":    {"time"},
		"4[value]":   {""},
		"5[name]":    {"sel_dept"},
		"5[value]":   {""},
		"6[name]":    {"sel_name"},
		"6[value]":   {""},
		"7[name]":    {"quick[]"},
		"7[value]":   {"N"},
		"8[name]":    {"quick[]"},
		"8[value]":   {"Y"},
		"9[name]":    {"quick[]"},
		"9[value]":   {"C"},
		"10[name]":   {"status[]"},
		"10[value]":  {"R0"},
		"11[name]":   {"status[]"},
		"11[value]":  {"S1"},
		"12[name]":   {"status[]"},
		"12[value]":  {"S3"},
		"13[name]":   {"status[]"},
		"13[value]":  {"S0"},
		"14[name]":   {"leave"},
		"14[value]":  {"A"},
		"15[name]":   {"brand[]"},
		"15[value]":  {"server"},
		"16[name]":   {"brand[]"},
		"16[value]":  {"cloud"},
		"17[name]":   {"brand[]"},
		"17[value]":  {"azure"},
		"18[name]":   {"brand[]"},
		"18[value]":  {"aws"},
		"19[name]":   {"search_mod"},
		"19[value]":  {"comp_name"},
		"20[name]":   {"search_val"},
		"20[value]":  {""},
		"21[name]":   {"perpage"},
		"21[value]":  {"100"},
		"hopeCal1":   {"0000-00-00"},
		"hopeCal2":   {"0000-00-00"},
		"Cal1":       {"0000-00-00"},
		"Cal2":       {"0000-00-00"},
		"quick[0]":   {"N"},
		"quick[1]":   {"Y"},
		"quick[2]":   {"C"},
		"status[0]":  {"R0"},
		"status[1]":  {"S1"},
		"status[2]":  {"S3"},
		"status[3]":  {"S0"},
		"leave":      {"A"},
		"brand[0]":   {"server"},
		"brand[1]":   {"cloud"},
		"brand[2]":   {"azure"},
		"brand[3]":   {"aws"},
		"search_mod": {"comp_name"},
		"perpage":    {"100"},
	}

	// 새로운 POST 요청 생성
	req, err := http.NewRequest("POST", config.GlobalConfig.Erp.WorkHistory.Url, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, fmt.Errorf("새 POST 요청 생성 실패: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 쿠키 포함하여 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP 요청 실패: %v", err)
	}
	defer resp.Body.Close()

	// 응답이 성공적이지 않은 경우
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// 응답 본문 파싱
	var result struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}

	if err := utils.ParseJSONBody(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("JSON 파싱 오류: %v", err)
	}

	// goquery를 사용해 HTML 문서를 파싱
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result.Msg))
	if err != nil {
		return nil, fmt.Errorf("HTML 문서 파싱 오류: %v", err)
	}

	var workHistories []WorkHistory

	doc.Find(".tbl_style3 tbody tr").Each(func(i int, row *goquery.Selection) {
		var workHistory WorkHistory

		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			columnIndex := j + 1
			columnContent := strings.TrimSpace(cell.Text())

			switch columnIndex {
			case 1:
				workHistory.Index = columnContent
			case 2:
				workHistory.ProcessingStatus = columnContent
			case 3:
				workHistory.RegistrationDate = columnContent
			case 4:
				workHistory.DesiredWorkTime = columnContent
			case 5:
				workHistory.ExpectedCompletionTime = columnContent
			case 6:
				workHistory.Worker = columnContent
			case 7:
				// <span> 태그의 클래스에 따라 내용 처리
				cell.Find("span").Each(func(k int, span *goquery.Selection) {
					class, _ := span.Attr("class")
					content := strings.TrimSpace(span.Text())

					switch class {
					case "cl_mediumblue":
						workHistory.RequestItemName = content
					case "tooltip-text":
						workHistory.RequestContent = content
					}
				})

				// <a> 태그의 링크 추출
				cell.Find("a").Each(func(k int, link *goquery.Selection) {
					class, _ := link.Attr("class")
					if class == "ln_inner" {
						href, exists := link.Attr("href")
						if exists {
							workHistory.RequestLink = href
						}
					}
				})
			case 8:
				workHistory.Subcategory = columnContent
			case 9:
				workHistory.CompanyName = columnContent
			case 10:
				workHistory.IP = columnContent
			case 11:
				workHistory.Brand = columnContent
			case 12:
				workHistory.RegistrationType = columnContent
			default:
				fmt.Printf("기타: %s\n", columnContent)
			}
		})

		workHistories = append(workHistories, workHistory)
	})

	return workHistories, nil
}
