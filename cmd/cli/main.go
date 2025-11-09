package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"salvadorsru/bob/internal/core/transpiler"
	"salvadorsru/bob/internal/lib/cli"
	"salvadorsru/bob/internal/lib/console"
	"salvadorsru/bob/internal/lib/failure"
	"salvadorsru/bob/internal/lib/file"
	"strings"
)

var version = "v0.1.0"

func collectFiles(input string, isFile, isFolder bool) (error, []string) {
	if isFile {
		return nil, []string{input}
	}
	if isFolder {
		return file.FindBobFiles(input)
	}
	return nil, nil
}

func panic(asJson bool, err *failure.Failure) {
	if err != nil {
		if asJson {
			data := map[string]map[string]string{
				"error": {
					"name":  err.Name,
					"value": err.Error(),
				},
			}
			jsonBytes, _ := json.MarshalIndent(data, "", "  ")
			console.Log(string(jsonBytes))
			os.Exit(1)
		}
		console.Panic(err.Error())
	}
}

func printResult(asJson bool, tables transpiler.TranspiledTable, actions transpiler.TranspiledActions) {
	if asJson {
		data := map[string]any{
			"tables":  tables.Get(),
			"actions": actions.Get(),
		}
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(asJson, failure.JsonParse)
		}
		console.Log(string(jsonBytes))
	} else {
		console.Success()
		console.Log(tables.ToString(), "\n\n", actions.ToString())
	}
}

func main() {
	console.Clear()
	argsErr, args := cli.ProcessArgs(version)
	if argsErr != nil {
		panic(args.AsJson, failure.Args)
	}

	if args.Driver == "" {
		console.Panic("driver not specified")
	}
	if args.OutputIsFile {
		console.Panic("output must be a folder")
	}

	driverErr, driver := transpiler.GetDriver(args.Driver)
	panic(args.AsJson, driverErr)

	if args.Query != "" {
		handleDirectQuery(*args, driver)
	} else {
		handleInputFiles(*args, driver)
	}
}

func handleDirectQuery(args cli.Args, driver transpiler.Driver) {
	transpileErr, tables, actions := transpiler.Transpile(driver, args.Query)
	panic(args.AsJson, transpileErr)

	if args.Output == "" {
		printResult(args.AsJson, *tables, *actions)
		return
	}

	files := []file.File{
		{Ref: "actions.sql", Content: actions.ToString()},
		{Ref: "tables.sql", Content: tables.ToString()},
	}
	file.WriteFiles(files, args.Output, args.OutputIsFolder)
	console.Success("transpiled to " + args.Output)
}

func handleInputFiles(args cli.Args, driver transpiler.Driver) {
	if args.Input == "" {
		console.Panic("invalid empty input")
	}

	err, filesList := collectFiles(args.Input, args.InputIsFile, args.InputIsFolder)
	if err != nil {
		panic(args.AsJson, failure.CollectFiles)
	}

	results := file.ReadFiles(filesList)
	var combinedInput strings.Builder
	var outputFiles []file.File

	for _, res := range results {
		if res.Err != nil {
			console.Log(res.Err)
			continue
		}

		actionErr, _, action := transpiler.Transpile(driver, res.Content)
		panic(args.AsJson, actionErr)

		if args.Output != "" {
			fileName := strings.TrimSuffix(filepath.Base(res.Ref), ".bob") + ".sql"
			outputFiles = append(outputFiles, file.File{Ref: fileName, Content: action.ToString()})
		}

		combinedInput.WriteString(res.Content)
		combinedInput.WriteByte('\n')
	}

	processCombined(args, driver, combinedInput.String(), outputFiles)
}

func processCombined(args cli.Args, driver transpiler.Driver, input string, files []file.File) {
	transpileErr, tables, actions := transpiler.Transpile(driver, input)
	panic(args.AsJson, transpileErr)

	if args.Output == "" {
		printResult(args.AsJson, *tables, *actions)
		return
	}

	files = append(files, file.File{Ref: "tables.sql", Content: tables.ToString()})
	file.WriteFiles(files, args.Output, args.OutputIsFolder)
	console.Success("transpiled to " + args.Output)
}
