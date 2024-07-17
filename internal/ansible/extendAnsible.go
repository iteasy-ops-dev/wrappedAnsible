package ansible

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/pkg/utils"
)

type extendAnsible struct {
	// DefaultAnsible DefaultAnsible
	// Public
	Ctx         context.Context
	Type        string   // Require: Name of playbook
	Name        string   // Require: Worker
	IPs         []string // Require: ips
	Account     string   // Require: account
	Password    string   // Require: remote server password
	Description string   // Description

	Options map[string]interface{}

	// Private
	inventory string // file.Name()
	playBook  string // file.Name()

	status   bool
	payload  string
	duration time.Duration
}

func (e *extendAnsible) generateInventoryPayload() []byte {
	var buffer bytes.Buffer
	for i := 0; i < len(e.IPs); i++ {
		entry := fmt.Sprintf(`%s ansible_user=%s ansible_password="%s"`+"\n", e.IPs[i], e.Account, e.Password)
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

func (e *extendAnsible) addStatus(
	status bool,
	stdoutStderr []byte,
	duration time.Duration,
) {
	e.status = status
	e.payload = string(stdoutStderr)
	e.duration = time.Duration(duration.Seconds())
}

func (e *extendAnsible) createExtraVars() string {
	// 꿀팁: https://docs.ansible.com/ansible/latest/playbook_guide/playbooks_variables.html
	jsonString, err := json.Marshal(e.Options)
	if err != nil {
		fmt.Println("JSON marshaling failed:", err)
	}
	return string(jsonString)
}

func (e *extendAnsible) excute() (*AnsibleProcessStatus, error) {
	status := true
	start := time.Now()
	e.createInventory()
	e.createPlaybook()
	defer utils.RemoveFile(e.inventory)

	fmt.Printf(
		"⚙️ Used Playbook: %s\n⚙️ Extra Vars: %s\n",
		e.playBook, e.createExtraVars(),
	)

	cmd := exec.CommandContext(
		e.Ctx,
		config.ANSIBLE_PLAYBOOK,
		e.playBook,
		config.OPTION_INVENTORY, e.inventory,
		config.OPTION_EXTRA_VARS, e.createExtraVars(),
	)

	// 채널을 생성하여 클라이언트 연결이 끊길 때 신호를 기다립니다.
	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-e.Ctx.Done():
			// 클라이언트 연결이 끊겼을 때 작업을 중단합니다.
			fmt.Println("❌ Client connection closed. Cancelling Ansible execution.")
			cmd.Process.Kill() // Ansible 프로세스를 강제로 종료합니다.
			status = false
		case <-done:
			return
		}
	}()

	// Ansible 작업을 수행합니다.
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(fmt.Errorf("❌ ERROR: stdoutStderr: %w", err))
		fmt.Println(fmt.Errorf("❌ ERROR: stdoutStderr: %s", stdoutStderr))
		status = false
	}

	duration := time.Since(start)

	e.addStatus(status, stdoutStderr, duration)

	s := newAnsibleProcessStatus(e)

	return s, nil
}
