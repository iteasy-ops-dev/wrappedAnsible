package ansible

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/pkg/utils"
)

type extendAnsible struct {
	// DefaultAnsible DefaultAnsible
	// Public
	Ctx            context.Context
	Type           string   // Require: Name of playbook
	Email          string   // Require: Worker
	Name           string   // Require: Worker
	IPs            []string // Require: ips
	Account        string   // Require: account
	Password       string   // Require: remote server password
	BecomePassword string   // Require: remote server root password
	Description    string   // Description

	Options map[string]interface{}

	// Private
	inventory string // file.Name()
	playBook  string // file.Name()

	status   bool
	payload  string
	duration time.Duration
}

// ansible_ssh_extra_args='-o HostKeyAlgorithms=+ssh-rsa' 추가.
// python 밑 os 낮은 버전들의 호환성을 위함.
func (e *extendAnsible) generateInventoryPayload() []byte {
	var buffer bytes.Buffer
	for i := 0; i < len(e.IPs); i++ {
		host := strings.Split(e.IPs[i], ":")
		if len(host) == 1 { // 기본포트
			entry := fmt.Sprintf(`%s ansible_host=%s ansible_user=%s ansible_password="%s" ansible_become_password="%s" ansible_ssh_extra_args='-o HostKeyAlgorithms=+ssh-rsa'`+"\n", host[0], host[0], e.Account, e.Password, e.BecomePassword)
			buffer.WriteString(entry)
			// fmt.Println(entry)
		} else { // 커스텀 포트
			entry := fmt.Sprintf(`%s ansible_host=%s ansible_port=%s ansible_user=%s ansible_password="%s" ansible_become_password="%s" ansible_ssh_extra_args='-o HostKeyAlgorithms=+ssh-rsa'`+"\n", host[0], host[0], host[1], e.Account, e.Password, e.BecomePassword)
			buffer.WriteString(entry)
			// fmt.Println(entry)
		}
	}
	return buffer.Bytes()
}

func (e *extendAnsible) createInventory() {
	r, err := utils.GenerateTempFile(e.generateInventoryPayload(), config.GlobalConfig.Ansible.Patterns.InventoryINI)
	if err != nil {
		log.Fatal(err)
	}
	e.inventory = r.Name()
}

func (e *extendAnsible) createPlaybook() {
	m := utils.GetFileListForMap(config.GlobalConfig.Ansible.PathStaticPlaybook)
	// m := utils.GetFileListForMap(config.PATH_STATIC_PLAYBOOK)
	path := fmt.Sprintf(`%s%s`, config.GlobalConfig.Ansible.PathStaticPlaybook, m[e.Type])
	// path := fmt.Sprintf(`%s%s`, config.PATH_STATIC_PLAYBOOK, m[e.Type])
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
		log.Println("JSON marshaling failed:", err)
	}
	return string(jsonString)
}

func (e *extendAnsible) excute() (*AnsibleProcessStatus, error) {
	status := true
	start := time.Now()
	e.createInventory()
	e.createPlaybook()
	// ansible 구동 후 임시 인벤토리 파일 삭제
	defer utils.RemoveFile(e.inventory)

	log.Printf("⚙️ Used Playbook: %s\n", e.playBook)
	log.Printf("⚙️ Extra Vars: %s\n", e.createExtraVars())

	cmd := exec.CommandContext(
		e.Ctx,
		config.GlobalConfig.Ansible.Playbook,
		e.playBook,
		config.GlobalConfig.Ansible.Options.Inventory, e.inventory,
		config.GlobalConfig.Ansible.Options.ExtraVars, e.createExtraVars(),
	)

	// 채널을 생성하여 클라이언트 연결이 끊길 때 신호를 기다립니다.
	done := make(chan struct{})
	defer close(done)

	// TODO: 강제로 종료되었을 경우 이후 어떻게 처리할지 정리
	// 예를 들어 DB에 따로 구분하여 정리한다던지
	go func() {
		select {
		case <-e.Ctx.Done():
			// 클라이언트 연결이 끊겼을 때 작업을 중단합니다.
			log.Println("❌ Client connection closed. Cancelling Ansible execution.")
			cmd.Process.Kill() // Ansible 프로세스를 강제로 종료합니다.
			status = false
		case <-done:
			return
		}
	}()

	// Ansible 작업을 수행합니다.
	stdoutStderr, err := cmd.CombinedOutput()

	if err != nil {
		log.Println(fmt.Errorf("❌ ERROR: stdoutStderr: %w", err))
		// log.Println(fmt.Errorf("❌ ERROR: stdoutStderr: %s", stdoutStderr))
		status = false
	}

	duration := time.Since(start)

	e.addStatus(status, stdoutStderr, duration)

	s := newAnsibleProcessStatus(e)

	return s, nil
}
