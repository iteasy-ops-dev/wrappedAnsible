package ansible

import (
	"context"
	"encoding/json"
	"log"
)

func Excuter(a iAnsible) ([]byte, error) {
	return a.excute()
}

func GetAnsibleFromFactory(ctx context.Context, jsonData []byte) iAnsible {
	e := extendAnsible{
		Ctx: ctx,
	}
	err := json.Unmarshal(jsonData, &e)
	// err := json.Unmarshal([]byte(jsonData), &e)
	if err != nil {
		log.Fatal(err)
	}

	return &e

	// if e.Type != "" {
	// 	return &e
	// } else {
	// 	return &e.DefaultAnsible
	// }
}
