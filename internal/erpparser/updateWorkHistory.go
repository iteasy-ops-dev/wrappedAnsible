package erpparser

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/internal/model"
)

// 날짜를 하루씩 증가시키는 함수
func GetNextDate(startDate string, daysToAdd int) string {
	layout := "2006-01-02"
	startTime, _ := time.Parse(layout, startDate)
	nextDate := startTime.AddDate(0, 0, daysToAdd)
	return nextDate.Format(layout)
}

// 날짜를 처리하는 함수
func ProcessDate(date string) error {
	// CookieJar 생성
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("쿠키 저장소 생성 실패:", err)
		return err
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
		fmt.Println("HTTP 요청 생성 실패:", err)
		return err
	}
	loginReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// HTTP 클라이언트 생성 (쿠키 포함)
	client := &http.Client{
		Jar: jar,
	}

	// 로그인 요청 보내기
	loginResp, err := client.Do(loginReq)
	if err != nil {
		fmt.Println("HTTP 요청 실패:", err)
		return err
	}
	defer loginResp.Body.Close()

	// 로그인 성공 여부 확인
	if loginResp.StatusCode != http.StatusOK {
		fmt.Println("로그인 실패: 상태 코드", loginResp.StatusCode)
		return err
	}

	payload := url.Values{
		"0[name]":    {"hopeCal1"},
		"0[value]":   {"0000-00-00"},
		"1[name]":    {"hopeCal2"},
		"1[value]":   {"0000-00-00"},
		"2[name]":    {"Cal1"},
		"2[value]":   {date},
		"3[name]":    {"Cal2"},
		"3[value]":   {date},
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
		"12[value]":  {"S2"},
		"13[name]":   {"status[]"},
		"13[value]":  {"S3"},
		"14[name]":   {"status[]"},
		"14[value]":  {"S4"},
		"15[name]":   {"leave"},
		"15[value]":  {"A"},
		"16[name]":   {"brand[]"},
		"16[value]":  {"server"},
		"17[name]":   {"brand[]"},
		"17[value]":  {"cloud"},
		"18[name]":   {"brand[]"},
		"18[value]":  {"azure"},
		"19[name]":   {"brand[]"},
		"19[value]":  {"aws"},
		"20[name]":   {"search_mod"},
		"20[value]":  {"comp_name"},
		"21[name]":   {"search_val"},
		"21[value]":  {""},
		"22[name]":   {"perpage"},
		"22[value]":  {"100"},
		"hopeCal1":   {"0000-00-00"},
		"hopeCal2":   {"0000-00-00"},
		"Cal1":       {date},
		"Cal2":       {date},
		"quick[0]":   {"N"},
		"quick[1]":   {"Y"},
		"quick[2]":   {"C"},
		"status[0]":  {"R0"},
		"status[1]":  {"S1"},
		"status[2]":  {"S2"},
		"status[3]":  {"S3"},
		"status[4]":  {"S4"},
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
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// 쿠키 포함하여 요청 보내기
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 응답이 성공적이지 않은 경우
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// 응답 본문 파싱
	var result struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
	}

	if err := parseJSONBody(resp.Body, &result); err != nil {
		return fmt.Errorf("JSON 파싱 오류: %v", err)
	}

	// Msg 필드 출력
	// fmt.Println("================================")
	// fmt.Printf("날짜: %s\n", date)

	// 추가적으로 goquery를 사용해 HTML 문서를 파싱할 수 있습니다.
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(result.Msg))
	if err != nil {
		return err
	}

	doc.Find(".tbl_style3 tbody tr").Each(func(i int, row *goquery.Selection) {
		m := model.NewWorkHistory()

		// fmt.Println("================================================")
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			columnIndex := j + 1
			columnContent := strings.TrimSpace(cell.Text())

			switch columnIndex {
			case 1:
				m.Setindex(columnContent)
				// fmt.Printf("인덱스: %s\n", columnContent)
			case 2:
				m.Setstatus(columnContent)
				// fmt.Printf("처리 상황: %s\n", columnContent)
			case 3:
				m.SetregistrationDate(columnContent)
				// fmt.Printf("등록일: %s\n", columnContent)
			case 4:
				m.SetdesiredWorkingHours(columnContent)
				// fmt.Printf("작업 희망시간: %s\n", columnContent)
			case 5:
				m.SetestimatedCompletionTime(columnContent)
				// fmt.Printf("완료 예정시간: %s\n", columnContent)
			case 6:
				m.Setworker(columnContent)
				// fmt.Printf("작업자: %s\n", columnContent)
			case 7:
				// a 태그에 있는 url은 추출하고 싶다.
				cell.Find("a").Each(func(k int, link *goquery.Selection) {
					url, exists := link.Attr("href")
					if exists {
						// fmt.Printf("링크 URL: %s\n", url)
						m.Seturl(url)
					}
				})
				if strings.Contains(columnContent, "    ") {
					splitIndex := strings.Index(columnContent, "    ")
					if splitIndex != -1 {
						item := strings.TrimSpace(columnContent[:splitIndex])
						content := strings.TrimSpace(columnContent[splitIndex+5:])
						// fmt.Printf("작업의뢰 항목: %s\n", item)
						// fmt.Printf("작업의뢰 내용: %s\n", content)
						m.SetworkRequestItems(item)
						if cell.Find("table").Length() > 0 {
							return // 내부 테이블이 있는 td를 건너뜁니다
						}
						m.SetworkRequestDetails(content)
					} else {
						// fmt.Printf("작업의뢰 항목: %s\n", columnContent)
						m.SetworkRequestItems(columnContent)
					}
				} else {
					m.SetworkRequestItems(columnContent)
					// fmt.Printf("작업의뢰 항목: %s\n", columnContent)
				}
			case 8:
				m.SetsubCategory(columnContent)
				// fmt.Printf("소분류: %s\n", columnContent)
			case 9:
				m.SetclientCompany(columnContent)
				// fmt.Printf("업체명: %s\n", columnContent)
			case 10:
				m.Setip(columnContent)
				// fmt.Printf("IP: %s\n", columnContent)
			case 11:
				m.Setbrand(columnContent)
				// fmt.Printf("브랜드: %s\n", columnContent)
			case 12:
				m.Setsection(columnContent)
				// fmt.Printf("등록구분: %s\n", columnContent)
			default:
				// fmt.Printf("기타: %s\n", columnContent)
			}
		})

		m.Put()
	})

	return nil
}

// parseJSONBody 함수: 응답 본문을 JSON으로 파싱
func parseJSONBody(body io.Reader, v interface{}) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return nil
}
