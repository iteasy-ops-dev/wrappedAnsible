package ansible

import (
	"encoding/json"
	"time"
)

type AnsibleProcessStatus struct {
	Type     string
	IPs      []string
	Name     string
	Account  string
	Status   bool
	Payload  string
	Duration time.Duration // int64

	Options map[string]interface{}
}

func newAnsibleProcessStatus(e *extendAnsible) *AnsibleProcessStatus {
	return &AnsibleProcessStatus{
		Type:     e.Type,
		IPs:      e.IPs,
		Name:     e.Name,
		Account:  e.Account,
		Status:   e.status,
		Payload:  e.payload,
		Duration: e.duration,
		Options:  e.Options,
	}
}

func (a *AnsibleProcessStatus) ToBytes() []byte {
	b, _ := json.Marshal(&a)
	return b
}
