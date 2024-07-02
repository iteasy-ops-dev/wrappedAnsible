package ansible

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/pkg/utils"
)

type extendAnsible struct {
	// DefaultAnsible DefaultAnsible
	Type        string   // Require: Name of playbook
	Name        string   // Require: Worker
	IPs         []string // Require: ips
	Account     string   // Require: account
	Password    string   // Require: remote server password
	Description string   // Description
	inventory   string   // file.Name()
	playBook    string   // file.Name()

	Options map[string]string
}

func (e *extendAnsible) generateInventoryPayload() []byte {
	var buffer bytes.Buffer
	for i := 0; i < len(e.IPs); i++ {
		entry := fmt.Sprintf(`%s ansible_user=%s ansible_password="%s" ansible_ssh_private_key_file=~/.ssh/control_node`+"\n", e.IPs[i], e.Account, e.Password)
		buffer.WriteString(entry)
	}
	return buffer.Bytes()
}

func (e *extendAnsible) createInventory() {
	r, err := utils.GenerateTempFile(e.generateInventoryPayload(), config.PATTERN_OF_INVENTORY_INI)
	if err != nil {
		log.Fatal(err)
	}
	e.inventory = r.Name()
}

func (e *extendAnsible) createPlaybook() {
	m := utils.GetFileListForMap(config.PATH_STATIC_PLAYBOOK)
	path := fmt.Sprintf(`%s%s`, config.PATH_STATIC_PLAYBOOK, m[e.Type])
	if utils.ExistFile(path) {
		e.playBook = path
	} else {
		log.Fatal("파일 없음.")
	}
}

func (e *extendAnsible) createExtraVars() string {
	var tmpString string
	for i, k := range e.Options {
		tmpString += fmt.Sprintf(`%s=%s `, i, k)
	}
	return tmpString
}

func (e *extendAnsible) excute() []byte {
	status := true
	e.createInventory()
	e.createPlaybook()
	fmt.Printf(
		"Name: %s\nAccount: %s\nInventory: %s\nPlaybook: %s\n",
		e.Name,
		e.Account,
		e.inventory,
		e.playBook,
	)
	fmt.Printf("Extra Vars: %s\n", e.createExtraVars())

	cmd := exec.Command(
		config.ANSIBLE_PLAYBOOK,
		e.playBook,
		config.OPTION_INVENTORY, e.inventory,
		config.OPTION_EXTRA_VARS, e.createExtraVars(),
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("ERROR: stdoutStderr: %s\n", err)
		status = false
	}

	type Res struct {
		Type    string
		Name    string
		Status  bool
		Payload string
	}

	s := Res{
		Type:    e.Type,
		Name:    e.Name,
		Status:  status,
		Payload: string(stdoutStderr),
	}

	// fmt.Printf("%s\n", stdoutStderr)
	// TODO: init 분기하지 말고 초기화 함수 따로 만들기
	// if e.Type != "init" {
	// 	o := Output(stdoutStderr).Debug(status)
	// 	b, _ := json.Marshal(o)
	// fmt.Printf("ID: %s\n", o.ID)
	// fmt.Printf("NAME: %s\n", o.Name)
	// for _, r := range o.Returns {
	// 	fmt.Printf("Name: %s\n", r.Name)
	// 	fmt.Printf("Host: %s\n", r.Host)
	// 	fmt.Printf("Action: %s\n", r.Action)
	// 	fmt.Printf("Msg: %s\n", r.Msg)
	// 	fmt.Printf("Failed: %t\n", r.Failed)
	// 	fmt.Printf("Results: %s\n", r.Results)
	// }
	// 	return b
	// }

	b, _ := json.Marshal(s)

	return b
}
