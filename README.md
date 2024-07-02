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