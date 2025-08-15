package cli

import (
	"errors"
	"os"
	"salvadorsru/bob/internal/core/console"
	"salvadorsru/bob/internal/core/utils"
)

type args struct {
	Driver         string
	InputIsFolder  bool
	InputIsFile    bool
	Input          string
	Query          string
	OutputIsFolder bool
	OutputIsFile   bool
	Output         string
}

func ProcessArgs(version string) (error, *args) {
	var args args

	for i, arg := range os.Args {
		if arg == "-v" || arg == "--version" {
			console.Log(version)
			os.Exit(0)
		}

		getFlagValue := func(i int, shortFlag string, longFlag string, isObligatory bool) (string, bool) {
			arg := os.Args[i]
			if (arg == shortFlag || (longFlag != "" && arg == longFlag)) && i+1 < len(os.Args) {
				nextArg := os.Args[i+1]
				if len(nextArg) > 0 && nextArg[0] == '-' {
					if isObligatory {
						console.Panic("expected value after %s, got flag %s\n", arg, nextArg)
						os.Exit(1)
					} else {
						return "", true
					}
				}
				return nextArg, true
			}
			return "", false
		}

		if val, ok := getFlagValue(i, "-d", "--driver", true); ok {
			args.Driver = val
		}

		if val, ok := getFlagValue(i, "-i", "--input", true); ok {
			pathType := utils.CheckPath(val)

			if utils.PathNotExist == pathType {
				return errors.New("invalid path type"), nil
			}

			if utils.PathIsDir == pathType {
				args.InputIsFolder = true
			} else if utils.PathIsFile == pathType {
				args.InputIsFile = true
			} else {
				return errors.New("invalid path type"), nil
			}

			args.Input = val
		}

		if val, ok := getFlagValue(i, "-q", "--query", true); ok {
			args.Query = val
		}

		if val, ok := getFlagValue(i, "-o", "--output", true); ok {
			pathType := utils.CheckPath(val)

			if utils.PathIsDir == pathType {
				args.OutputIsFolder = true
			} else if utils.PathIsFile == pathType {
				args.OutputIsFile = true
			} else {
				return errors.New("invalid path type"), nil
			}

			args.Output = val
		}
	}

	return nil, &args
}
