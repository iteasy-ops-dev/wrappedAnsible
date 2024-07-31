# ITEASY Wrapped Ansible

### TODO
- [ ] for ansible galaxy
- [ ] model
- [ ] db
- [ ] log

### ADD
- url은 handler/handlers.go 에 추가한다.
  - handler는 handler 폴더에 추가 한다.
- ansible 기능은 ansible 폴더에 추가한다.
  - ansible 기능 추가시 requiments.yml에 git 주소 추가한다.
  
### 대략 적인 Flow
- ansible roles는 따로 개발되어 회사 git 계정에서 가져온다.
- 해당 roles을 호출하여 실행하는 간략한 playbook이 ansible에 있다.
- url을 통해 해당 타입의 playbook을 실행한다.

### config
```json
{
	"default": {
		"host": "{{ Backend Host }}"
	},
	"jwt": {
		"key": "",
		"token_name": ""
	},
	"ansible": {
		"playbook": "ansible-playbook",
		"options": {
			"extra_vars": "--extra-vars",
			"inventory": "--inventory"
		},
		"default_extra_vars": "",
		"patterns": {
			"inventory_ini": "",
			"ansible_yml": ""
		},
		"path_static_playbook": ""
	},
	"mongodb": {
		"url": "",
		"database": "",
		"collections": {
			"ansible_process_status": "",
			"auth": ""
		}
	},
	"erp": {
		"base_url": "",
		"login": {
			"url": "",
			"admin_id": "",
			"admin_passwd": "",
			"allow_type": "",
			"login_btn": ""
		}
	},
	"smtp": {
		"host": "",
		"port": "",
		"from": ""
	}
}
```