package arguments

import "flag"

type Arguments struct {
	ParametersPath string
	Debug          bool
	Migrate        bool
}

func GetArguments() Arguments {
	var args Arguments

	flag.StringVar(&args.ParametersPath, "parameters", "./parameters.yaml", "Path to parameters file")
	flag.BoolVar(&args.Debug, "debug", false, "Debug mode")
	flag.BoolVar(&args.Migrate, "migrate", false, "Migrate database")
	flag.Parse()

	return args
}
