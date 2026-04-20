package settings

import (
	"encoding/json"

	"github.com/sater-151/todo-list/internal/init/arguments"
	"github.com/sater-151/todo-list/internal/init/config"
	"github.com/sater-151/todo-list/internal/init/parameters"
	"github.com/sater-151/todo-list/pkg/utils"
)

type AppInit struct {
	Arguments  arguments.Arguments
	Config     config.Config
	Parameters parameters.Parameters
}

func Init() (app AppInit, err error) {
	defer utils.AddFuncLabel("[init]", err)

	args := arguments.GetArguments()

	cfg, err := config.GetConfig()
	if err != nil {
		return
	}

	params, err := parameters.GetParameters(args.ParametersPath)
	if err != nil {
		return
	}

	if args.Debug {
		argsJson, _ := json.MarshalIndent(args, "", "  ")
		cfgJson, _ := json.MarshalIndent(cfg, "", "  ")
		paramsJson, _ := json.MarshalIndent(params, "", "  ")

		println("Arguments:\n", string(argsJson))
		println("Config:\n", string(cfgJson))
		println("Parameters:\n", string(paramsJson))
	}

	return AppInit{
		Arguments:  args,
		Config:     cfg,
		Parameters: params,
	}, nil

}
