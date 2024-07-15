package ansible

type iAnsible interface {
	createInventory()
	createPlaybook()
	excute() (*AnsibleProcessStatus, error)
}
