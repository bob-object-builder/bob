package main

import (
	"fmt"
	"os"
	"salvadorsru/bob/internal/core/console"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/file"
	"salvadorsru/bob/internal/core/response"
	"salvadorsru/bob/internal/core/utils"
	"salvadorsru/bob/internal/transpiler"
	"sync"
)

var version string = "0.0.0"

func main() {
	console.Clear()

	var queryFile string
	var driverName string
	var queryString string
	var outputFile string
	var searchMode bool
	var queries string
	var searchTarget string = "."

	var tables string
	var actions string
	var transpilerError error

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

		if val, ok := getFlagValue(i, "-i", "--input", true); ok {
			queryFile = val
		}

		if val, ok := getFlagValue(i, "-d", "--driver", true); ok {
			driverName = val
		}

		if val, ok := getFlagValue(i, "-q", "--query", true); ok {
			queryString = val
		}

		if val, ok := getFlagValue(i, "-o", "--output", true); ok {
			outputFile = val
		}

		if arg == "-s" || arg == "--search" {
			searchMode = true
			if val, hasValue := getFlagValue(i, "-s", "--search", false); hasValue {
				if val != "" {
					searchTarget = val
				}
			}
		}
	}

	if driverName == "" {
		console.Panic("no driver specified. Use -d <driver> or --driver <driver> (mariadb, postgresql, sqlite)")
	}

	if searchMode {
		files, err := utils.FindBobFiles(searchTarget)

		if err != nil {
			console.Panic("searching for .bob files:", err)
		}

		if len(files) == 0 {
			console.Panic("no .bob files found in the current directory and subdirectories.")
		}

		var wg sync.WaitGroup

		results := make([]file.File, len(files))

		for i, fileContent := range files {
			wg.Add(1)
			go func(idx int, filename string) {
				defer wg.Done()
				queryBytes, err := os.ReadFile(filename)
				if err != nil {
					results[idx] = file.File{Content: "", Err: response.Error("reading %s", filename), Ref: filename}
					return
				}
				results[idx] = file.File{Content: string(queryBytes), Err: nil, Ref: filename}
			}(i, fileContent)
		}
		wg.Wait()

		var allInput string
		for _, res := range results {
			if res.Err != nil {
				console.Log(res.Err)
				continue
			}
			allInput += res.Content + "\n"
		}

		transpilerError, tables, actions = transpiler.Transpile(drivers.Motor(driverName), allInput)
	} else {
		var input string

		if queryString != "" {
			input = queryString
		} else {
			if queryFile == "" {
				console.Panic("no input file specified. Use -i <file> or provide a query with -q <query>")
			}
			queryBytes, err := os.ReadFile(queryFile)
			if err != nil {
				console.Panic("not found", queryFile)
			}
			input = string(queryBytes)
		}

		transpilerError, tables, actions = transpiler.Transpile(drivers.Motor(driverName), input)
	}

	if transpilerError != nil {
		console.Panic(transpilerError.Error())
	}

	queries = fmt.Sprintf("%s\n\n%s\n", tables, actions)
	if outputFile != "" {
		file, err := os.Create(outputFile)
		if err != nil {
			console.Panic("creating file", err)
		}
		_, err = file.Write([]byte(queries))
		if err != nil {
			file.Close()
			console.Panic("writing to file", err)
		}
		console.Success("file created at", outputFile)
		defer file.Close()
	} else {
		console.Log(queries)
	}
}
