package handlers

import (
	"net/http"

	"iteasy.wrappedAnsible/internal/erpparser"
	"iteasy.wrappedAnsible/pkg/utils"
)

type erpParserReq struct {
	URL string `json:"url"` // JSON의 "url" 필드를 매핑
}

func erpParser(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodPost); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}
	data, err := utils.ParseRequestBody[erpParserReq](r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := data.URL

	if url == "" {
		http.Error(w, "URL not provided", http.StatusBadRequest)
		return
	}

	e := erpparser.New(url)
	e.Parsing()

	_httpResponse(w, http.StatusOK, e)
}
