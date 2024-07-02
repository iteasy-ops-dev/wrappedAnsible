package ansible

import (
	"encoding/json"
	"log"
)

func Excuter(a iAnsible) []byte {
	return a.excute()
}

func GetAnsibleFromFactory(jsonData []byte) iAnsible {
	e := extendAnsible{}
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
