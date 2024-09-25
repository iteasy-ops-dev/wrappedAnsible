package erpparser

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func WorkHistoryParsing() ([]WorkHistory, error) {

	client, _ := newLoginClient()
	payload := createPayload("")
	doc, _ := newDoc(client, payload)

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
				// log.Printf("기타: %s\n", columnContent)
			}
		})

		workHistories = append(workHistories, workHistory)
	})

	return workHistories, nil
}
