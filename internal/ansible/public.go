package ansible

import (
	"context"
	"errors"
	"net/http"
)

type GennerateInitType struct {
	Ctx      context.Context
	JsonData []byte
}

type GennerateHttpRequestType struct {
	Ctx context.Context
	R   *http.Request
}

func Excuter(a iAnsible) ([]byte, error) {
	return a.excute()
}

// func GetAnsibleFromFactory(ctx context.Context, jsonData []byte) iAnsible {
// 	e := extendAnsible{
// 		Ctx: ctx,
// 	}
// 	err := json.Unmarshal(jsonData, &e)
// 	// err := json.Unmarshal([]byte(jsonData), &e)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

//		return &e
//	}

func GetAnsibleFromFactory(v interface{}) (iAnsible, error) {
	switch v := v.(type) {
	case GennerateInitType:
		return generateInitAnsible(v)
	case GennerateHttpRequestType:
		return generateHttpAnsible(v)
	default:
		return nil, errors.New("구성할 수 없는 구조체 타입")
	}
}
