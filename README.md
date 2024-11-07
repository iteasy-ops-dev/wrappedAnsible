# ITEASY Wrapped Ansible

TODO
---------
- 해야 할일은 TODO: 를 검색해서 빠르게 해결하자

Infomation
---------
- ITEASY 운영 플랫폼에서 백엔드를 담당합니다.

실행
---------
```sh
# Common
docker-compose -f docker-compose-init.yml -p mongo up -d
# DEV
docker-compose -f docker-compose-dev.yml -p backend up -d
# PROD
docker-compose -f docker-compose-prod.yml -p backend up -d
```

정지
---------
```sh
# Common
docker-compose -f docker-compose-init.yml -p mongo down
# DEV
docker-compose -f docker-compose-dev.yml -p backend down --rmi all
# PROD
docker-compose -f docker-compose-prod.yml -p backend down --rmi all
```

Function
---------
- 만들어진 ansible-playbook을 구동합니다.
- ansible을 사용하기 위해 erp를 파싱합니다.
- 회원 관리를 담당합니다.

주의 사항
---------
- 사용하는 버전은 다음과 같습니다.
  - OS: ubuntu:20.04
  - Ansible: 2.13.13
  - Jinja: 3.1.4

업데이트 방법
---------
- 플레이북 추가 및 수정
  - 수정
    - 기존의 ansible playbook이 존재하는 경우
      - 수정된 플레이북을 github에 main으로 push
      - backend 컨테이너 재시작
  - 추가
      - github에 main으로 push 되어 있어야 하고,
      - 최상위 ansible 폴더에 플레이북 추가
      - ```yaml
          # 예시
        - hosts: all
          become: true
          roles:
            - ansible.roles.webhost_manager
        ```
      - 같은 폴더내의 requirements.yml에 github 정보 추가
      - ```yml
        # 예시
        - src: https://github.com/iteasy-ops-dev/ansible.roles.webhost_manager.git
          scm: git
          version: main
        ```
      - 재 배포

폴더 구조
---------
```
wrappedAnsible/
├── .vscode/
├── ansible/          # ansible 플레이북
├── cmd/
│   └── main.js       # 진입점
├── configs/
│   └── config.js     # 전역 설정
├── data/             # DB 데이터 마운트 폴더
├── internal/
│   ├── ansible/      # ansible 구동
│   ├── erpparser/    # Erp 파싱
│   ├── handler/      # 라우터 핸들러
│   ├── model/        # DB
│   └── router/       # 라우터
└── pkg/
    └── utils/        # 전역적으로 사용할 유틸리티
```

config
---------
```json
{
	"default": {
		"host": "{{ Backend Host }}",
    "admin":"{{ admin }}",
		"password": "{{ admin init password }}"
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

HISTORY
---------
`24/11/07`
- [ ] root 패스워드 관련 분기가 필요
  - [ ] 해당 ansible에 관련하여 become 메소드를 su 로 변경할 경우
  - [ ] map 형식으로 해당 타입의 메소드 여부를 확인
  - [ ] util에 해당 함수 추가

`24/09/20`
- 이진주 매니저님
  - [x] erp 파싱 문제(https://admin.ksidc.net/service/request_info/?mem_idx=40396&idx=229099)
    - [x] text로 파싱하던 것을 html로 파싱하는 것으로 해결
    - [x] 정규 표현식 수정
    - [x] ID, PW는 같이 파싱하여 배열의 인덱스로 추출
  - [ ] 원격 서버 파이썬 버전 낮은 문제 발생(python 2.7 미만 버전)
    - [ ] ansible raw 모듈에 대한 고찰이 필요. 아니면 쉘 스크립트로?
  - [ ] vhost 복수개 입력 불가능 문제

`24/08/26`
- 서강호 매니저님
  - [x] vhost의 경로변경