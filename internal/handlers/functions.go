package handlers

import (
	"net/http"

	config "iteasy.wrappedAnsible/configs"
	"iteasy.wrappedAnsible/pkg/utils"
)

func functions(w http.ResponseWriter, r *http.Request) {
	if err := _allowMethod(w, r, http.MethodGet); err != nil {
		return
	}
	if err := _validateToken(w, r); err != nil {
		return
	}

	result := make([]string, 0)
	l := utils.GetFileList(config.GlobalConfig.Ansible.PathStaticPlaybook)
	// l := utils.GetFileList(config.PATH_STATIC_PLAYBOOK)
	for _, name := range l {
		if utils.CheckExtension(name, `.yml`) {
			// 초기화에 필요한 yml이므로 제외
			if name == "requirements.yml" || name == "init.yml" {
				continue
			}
			result = append(result, utils.TruncationExtension(name))
		}
	}

	_httpResponse(w, http.StatusOK, result)
}
