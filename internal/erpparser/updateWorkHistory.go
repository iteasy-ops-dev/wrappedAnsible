package erpparser

import (
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

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

	client, _ := newLoginClient()
	payload := createPayload(date)
	doc, _ := newDoc(client, payload)

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
				// erp 링크
				cell.Find("a").Each(func(k int, link *goquery.Selection) {
					class, _ := link.Attr("class")
					if class == "ln_inner" {
						href, exists := link.Attr("href")
						if exists {
							m.Seturl(href)
						}
					}
				})

				// <span> 태그의 클래스에 따라 내용 처리
				cell.Find("span").Each(func(k int, span *goquery.Selection) {
					class, _ := span.Attr("class")
					content := strings.TrimSpace(span.Text())

					switch class {
					case "cl_mediumblue": // 작업의뢰 항목
						m.SetworkRequestItems(content)
					case "tooltip-text": // 세부내용
						m.SetworkRequestDetails(content)
					}
				})
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
				fmt.Printf("기타: %s\n", columnContent)
			}
		})

		m.Put()
	})

	return nil
}
