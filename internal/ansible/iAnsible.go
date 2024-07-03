package ansible

type iAnsible interface {
	createInventory()
	createPlaybook()
	excute() ([]byte, error)
}
