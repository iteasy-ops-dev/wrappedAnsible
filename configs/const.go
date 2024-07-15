package config

const ANSIBLE_PLAYBOOK = "ansible-playbook"

const OPTION_EXTRA_VARS = "--extra-vars"
const OPTION_INVENTORY = "--inventory"

const ANSIBLE_DEFAULT_EXTRA_VARS = "ansible_user=root ansible_password=1123"

const PATTERN_OF_INVENTORY_INI = "inventory-*.ini"
const PATTERN_OF_ANSIBLE_YML = "ansible-*.yml"

// Create Temp File Flag
// const FLAG_INVENTORY = "inventory"
// const FLAG_PLAYBOOK = "playbook"

const PATH_STATIC_PLAYBOOK = "ansible/"

// DB
const MONGODB_URL = "mongodb://mongo:27017"
const MONGODB_DATABASE = "iteasy-ops-dev"
const COLLECTION_ANSIBLE_PROCESS_STATUS = "ansible_process_status"
