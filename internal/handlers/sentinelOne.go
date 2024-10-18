// TODO: 센티널원에서 가져오는 데이터를 가지고
// 어떤 유효한 정보로 가공할 것인지에 대한 아이디어 필요

package handlers

import (
	"fmt"
	"net/http"

	"iteasy.wrappedAnsible/internal/sentinelone"
	"iteasy.wrappedAnsible/pkg/utils"
)

type SentinelOneAPIReq struct {
	ApiKey string `json:"apikey"`
	Type   string `json:"type"`
}

func updateSentinelOne(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	b, err := utils.ParseRequestBody[SentinelOneAPIReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s, err := sentinelone.GenerateSentinelOne(b.Type, b.ApiKey)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusNotFound)
		return
	}

	if err := s.Update(); err != nil {
		http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
		return
	}

	_httpResponse(w, http.StatusOK, nil)
}
