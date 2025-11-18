package main

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/core/transpiler"
	"salvadorsru/bob/internal/lib/cli"
	"salvadorsru/bob/internal/lib/console"
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

func panic(asJson bool, asDaemon bool, err *failure.Failure) {
	if err == nil {
		return
	}

	if asJson {
		data := map[string]map[string]string{
			"error": {
				"type":    err.Name,
				"message": strings.TrimSpace(err.Error()),
			},
		}
		jsonBytes, _ := json.MarshalIndent(data, "", "  ")
		console.Log(string(jsonBytes))
	} else {
		console.Panic(err.Error())
	}

	if asDaemon {
		console.Log("__END__")
		return
	}

	os.Exit(1)
}

func printResult(asJson bool, asDaemon bool, tables *transpiler.TranspiledTable, actions *transpiler.TranspiledActions) {
	if asJson {
		data := map[string]any{}
		if tables != nil {
			data["tables"] = tables.Get()
		}
		if actions != nil {
			data["actions"] = actions.Get()
		}
		jsonBytes, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			panic(asJson, asDaemon, failure.JsonParse)
		}
		console.Log(string(jsonBytes))
		return
	}

	console.Success()
	if tables != nil {
		console.Log(tables.ToString())
	}
	if actions != nil {
		console.Log()
		console.Log(actions.ToString())
	}
}

func main() {
	console.Clear()
	argsErr, args := cli.ProcessArgs(version)
	if argsErr != nil {
		panic(args.AsJson, args.AsDaemon, failure.MalformedArgs)
	}

	driverErr, driver := transpiler.GetDriver(args.Driver)
	panic(args.AsJson, args.AsDaemon, driverErr)

	if args.AsDaemon {
		handleDaemonQuery(*args, driver)
	} else if args.Query != "" {
		handleDirectQuery(*args, driver)
	} else {
		handleInputFiles(*args, driver)
	}
}

func handleDaemonQuery(args cli.Args, driver transpiler.Driver) {
	scanner := bufio.NewScanner(os.Stdin)
	var queryBuilder strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		if line == "__EXIT__" {
			os.Exit(0)
		}

		if line == "__END__" {
			query := strings.TrimSpace(queryBuilder.String())
			if query != "" {
				transpileErr, tables, actions := transpiler.Transpile(driver, query)
				panic(args.AsJson, args.AsDaemon, transpileErr)
				printResult(args.AsJson, args.AsDaemon, tables, actions)
			}
			console.Log("__END__")
			queryBuilder.Reset()
			continue
		}

		queryBuilder.WriteString(line)
		queryBuilder.WriteByte('\n')
	}

	if err := scanner.Err(); err != nil {
		panic(args.AsJson, args.AsDaemon, failure.IO)
	}
}

func handleDirectQuery(args cli.Args, driver transpiler.Driver) {
	transpileErr, tables, actions := transpiler.Transpile(driver, args.Query)
	panic(args.AsJson, args.AsDaemon, transpileErr)

	if args.Output == "" {
		printResult(args.AsJson, args.AsDaemon, tables, actions)
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
		panic(args.AsJson, args.AsDaemon, failure.InvalidInput)
	}

	err, filesList := collectFiles(args.Input, args.InputIsFile, args.InputIsFolder)
	if err != nil {
		panic(args.AsJson, args.AsDaemon, failure.CollectFiles)
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
		panic(args.AsJson, args.AsDaemon, actionErr)

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
	panic(args.AsJson, args.AsDaemon, transpileErr)

	if args.Output == "" {
		printResult(args.AsJson, args.AsDaemon, tables, actions)
		return
	}

	files = append(files, file.File{Ref: "tables.sql", Content: tables.ToString()})
	file.WriteFiles(files, args.Output, args.OutputIsFolder)
	console.Success("transpiled to " + args.Output)
}
