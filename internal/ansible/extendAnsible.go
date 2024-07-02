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
		"⚙️ Used Playbook: %s\n⚙️ Extra Vars: %s\n",
		e.playBook, e.createExtraVars(),
	)

	cmd := exec.Command(
		config.ANSIBLE_PLAYBOOK,
		e.playBook,
		config.OPTION_INVENTORY, e.inventory,
		config.OPTION_EXTRA_VARS, e.createExtraVars(),
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ ERROR: stdoutStderr: %s\n", err)
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

	b, _ := json.Marshal(s)

	return b
}