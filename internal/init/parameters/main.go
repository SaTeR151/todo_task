package parameters

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sater-151/todo-list/pkg/utils"
)

type Parameters struct {
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func GetParameters(parametersPath string) (params Parameters, err error) {
	defer utils.AddFuncLabel("[init-get-parameters]", err)

	if err = cleanenv.ReadConfig(parametersPath, &params); err != nil {
		return
	}

	return
}
