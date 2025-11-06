package main

import (
	"path/filepath"
	"salvadorsru/bob/internal/core/transpiler"
	"salvadorsru/bob/internal/lib/cli"
	"salvadorsru/bob/internal/lib/console"
	"salvadorsru/bob/internal/lib/file"
	"strings"
)

var version = "v0.1.0"

func collectFiles(input string, isFile, isFolder bool) ([]string, error) {
	if isFile {
		return []string{input}, nil
	}
	if isFolder {
		return file.FindBobFiles(input)
	}
	return nil, nil
}

func mustPanicOnError(err error) {
	if err != nil {
		console.Panic(err.Error())
	}
}

func main() {
	argsErr, args := cli.ProcessArgs(version)
	mustPanicOnError(argsErr)

	if args.Driver == "" {
		console.Panic("driver not specified")
	}
	if args.OutputIsFile {
		console.Panic("output must be a folder")
	}

	driverErr, driver := transpiler.GetDriver(args.Driver)
	mustPanicOnError(driverErr)

	if args.Query != "" {
		handleDirectQuery(*args, driver)
	} else {
		handleInputFiles(*args, driver)
	}
}

func handleDirectQuery(args cli.Args, driver transpiler.Driver) {
	transpileErr, tables, actions := transpiler.Transpile(driver, args.Query)
	mustPanicOnError(transpileErr)

	if args.Output == "" {
		console.Success()
		console.Log(tables, actions)
		return
	}

	file.WriteFiles([]file.File{
		{Ref: "actions.sql", Content: actions},
		{Ref: "tables.sql", Content: tables},
	}, args.Output, args.OutputIsFolder)
}

func handleInputFiles(args cli.Args, driver transpiler.Driver) {
	if args.Input == "" {
		console.Panic("invalid empty input")
	}

	files, err := collectFiles(args.Input, args.InputIsFile, args.InputIsFolder)
	mustPanicOnError(err)

	results := file.ReadFiles(files)

	var (
		combinedInput strings.Builder
		outputFiles   []file.File
	)

	for _, res := range results {
		if res.Err != nil {
			console.Log(res.Err)
			continue
		}

		actionErr, _, action := transpiler.Transpile(driver, res.Content)
		mustPanicOnError(actionErr)

		if args.Output != "" {
			fileName := strings.TrimSuffix(filepath.Base(res.Ref), ".bob") + ".sql"
			outputFiles = append(outputFiles, file.File{Ref: fileName, Content: action})
		}

		combinedInput.WriteString(res.Content)
		combinedInput.WriteByte('\n')
	}

	processCombined(args, driver, combinedInput.String(), outputFiles)
}

func processCombined(args cli.Args, driver transpiler.Driver, input string, files []file.File) {
	tablesErr, tables, actions := transpiler.Transpile(driver, input)
	mustPanicOnError(tablesErr)

	if args.Output == "" {
		console.Clear()
		console.Success()
		console.Log(tables, actions)
		return
	}

	files = append(files, file.File{Ref: "tables.sql", Content: tables})
	file.WriteFiles(files, args.Output, args.OutputIsFolder)
	console.Success("transpiled to " + args.Output)
}
