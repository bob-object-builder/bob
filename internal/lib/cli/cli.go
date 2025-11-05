package cli

import (
	"errors"
	"os"
	"salvadorsru/bob/internal/lib/checker"
	"salvadorsru/bob/internal/lib/console"
)

type Args struct {
	Driver         string
	Query          string
	InputIsFolder  bool
	InputIsFile    bool
	Input          string
	OutputIsFolder bool
	OutputIsFile   bool
	Output         string
}

func ProcessArgs(version string) (error, *Args) {
	var args Args

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
			pathType := checker.ValidatePath(val)

			if checker.PathNotExist == pathType {
				return errors.New("invalid path type"), nil
			}

			if checker.PathIsDir == pathType {
				args.InputIsFolder = true
			} else if checker.PathIsFile == pathType {
				args.InputIsFile = true
			} else {
				return errors.New("invalid path type"), nil
			}

			if val == "" {
				return errors.New("empty query"), nil
			}

			args.Input = val
		}

		if val, ok := getFlagValue(i, "-q", "--query", true); ok {
			args.Query = val
		}

		if val, ok := getFlagValue(i, "-o", "--output", true); ok {
			pathType := checker.CheckPath(val)

			if checker.PathIsDir == pathType {
				args.OutputIsFolder = true
			} else if checker.PathIsFile == pathType {
				args.OutputIsFile = true
			} else {
				return errors.New("invalid path type"), nil
			}

			args.Output = val
		}
	}

	return nil, &args
}
