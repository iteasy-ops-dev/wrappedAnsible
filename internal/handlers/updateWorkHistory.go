// package handlers

// import (
// 	"io"
// 	"log"
// 	"net/http"
// 	"sync"

// 	"iteasy.wrappedAnsible/internal/erpparser"
// 	"iteasy.wrappedAnsible/pkg/utils"
// )

// func UpdateWorkHistory(w http.ResponseWriter, r *http.Request) {
// 	if err := _allowMethod(w, r, http.MethodPost); err != nil {
// 		return
// 	}
// 	if err := _validateToken(w, r); err != nil {
// 		return
// 	}
// 	type req struct {
// 		StartDate string `json:"start_date"`
// 		NumDays   int    `json:"num_days"`
// 	}

// 	// 요청 바디에서 JSON 데이터를 파싱합니다.
// 	data, err := utils.ParseRequestBody[req](r)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	var wg sync.WaitGroup
// 	for i := 0; i < data.NumDays; i++ {
// 		date := erpparser.GetNextDate(data.StartDate, i)
// 		wg.Add(1) // WaitGroup 카운트 증가
// 		go func(date string) {
// 			defer wg.Done() // 고루틴 완료 후 WaitGroup 카운트 감소
// 			if err := erpparser.ProcessDate(date); err != nil {
// 				log.Printf("날짜 %s 처리 중 오류 발생: %v", date, err)
// 			}
// 		}(date)
// 	}
// 	wg.Wait() // 모든 고루틴이 완료될 때까지 대기

// 	io.WriteString(w, "update data!\n")
// }

package handlers

import (
	"log"
	"net/http"

	"iteasy.wrappedAnsible/internal/erpparser"
	"iteasy.wrappedAnsible/pkg/utils"
)

func updateWorkHistory(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	type req struct {
		StartDate string `json:"start_date"`
		NumDays   int    `json:"num_days"`
	}

	// 요청 바디에서 JSON 데이터를 파싱합니다.
	data, err := utils.ParseRequestBody[req](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 즉시 응답을 반환
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusOK)

	_httpResponse(w, http.StatusOK, nil)

	// 백그라운드에서 고루틴 실행
	go func() {
		for i := 0; i < data.NumDays; i++ {
			date := erpparser.GetNextDate(data.StartDate, i)
			go func(date string) {
				if err := erpparser.ProcessDate(date); err != nil {
					log.Printf("날짜 %s 처리 중 오류 발생: %v", date, err)
				}
			}(date)
		}
	}()
}
