package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"iteasy.wrappedAnsible/internal/ansible"
	"iteasy.wrappedAnsible/internal/handlers"
	"iteasy.wrappedAnsible/internal/router"
)

func init() {

	fmt.Println("⚙️ Wrapped Ansible Server Init.")
	initJsonData := `{
			"type": "init",
			"name": "서버 초기화 실행.",
		  "options": {}
	  }
	`
	extendAnsible := ansible.GetAnsibleFromFactory([]byte(initJsonData))
	r := ansible.Excuter(extendAnsible)

	type InitReturn struct {
		Status bool
	}
	var i InitReturn
	err := json.Unmarshal(r, &i)
	if err != nil || !i.Status {
		panic("❌ 서버 초기화 실패.")
	}
}

func main() {
	fmt.Println("✅ Welcome Wrapped Ansible Server. PORT: 8080")

	log.Fatal(http.ListenAndServe(":8080", handlers.CorsMiddleware(router.NewRouter())))
}
