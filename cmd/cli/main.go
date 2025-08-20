package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"salvadorsru/bob/internal/core/cli"
	"salvadorsru/bob/internal/core/console"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/file"
	"salvadorsru/bob/internal/core/response"
	"salvadorsru/bob/internal/core/utils"
	"salvadorsru/bob/internal/transpiler"
)

var version = "v0.0.0"

func main() {
	argsError, args := cli.ProcessArgs(version)
	if argsError != nil {
		console.Panic(argsError.Error())
		return
	}

	if args.Driver == "" {
		console.Panic("driver not specified")
		return
	}
	if args.OutputIsFile {
		console.Panic("output must be a folder")
		return
	}

	var (
		allInput      strings.Builder
		filesToCreate []file.File
		hasOutput     = args.Output != ""
	)

	if args.Query != "" {

		transpileError, tables, actions := transpiler.Transpile(drivers.Motor(args.Driver), args.Query)
		if transpileError != nil {
			console.Panic(transpileError.Error())
			return
		}

		if !hasOutput {
			console.Success()
			console.Log(tables, actions)
		} else {
			writeFiles([]file.File{
				{
					Ref:     "actions.sql",
					Content: actions,
				},
				{
					Ref:     "tables.sql",
					Content: tables,
				},
			}, args.Output, args.OutputIsFolder)
		}

	} else {
		fileList := collectFiles(args.Input, args.InputIsFile, args.InputIsFolder)
		results := readFiles(fileList)

		for _, result := range results {
			if result.Err != nil {
				console.Log(result.Err)
				continue
			}

			actionErr, _, action := transpiler.TranspileActions(drivers.Motor(args.Driver), result.Content)
			if actionErr != nil {
				console.Panic(actionErr.Error())
				return
			}

			if hasOutput {
				fileName := strings.TrimSuffix(filepath.Base(result.Ref), ".bob") + ".sql"
				filesToCreate = append(filesToCreate, file.File{
					Ref:     fileName,
					Content: action,
				})
			}

			allInput.WriteString(result.Content)
			allInput.WriteByte('\n')
		}

		if !hasOutput {

			tablesErr, tables, actions := transpiler.Transpile(drivers.Motor(args.Driver), allInput.String())
			if tablesErr != nil {
				console.Panic(tablesErr.Error())
				return
			}

			console.Success()
			console.Log(tables, actions)

		} else {

			tablesErr, tables, _ := transpiler.TranspileTables(drivers.Motor(args.Driver), allInput.String())
			if tablesErr != nil {
				console.Panic(tablesErr.Error())
				return
			}

			filesToCreate = append(filesToCreate, file.File{
				Ref:     "tables.sql",
				Content: tables,
			})

			writeFiles(filesToCreate, args.Output, args.OutputIsFolder)
			console.Success("transpiled to " + args.Output)
		}
	}
}

func collectFiles(input string, isFile, isFolder bool) []string {
	var fileList []string
	if isFile {
		return append(fileList, input)
	}
	if isFolder {
		files, err := utils.FindBobFiles(input)
		if err != nil {
			console.Panic(err.Error())
			return nil
		}
		return append(fileList, files...)
	}
	return fileList
}

func readFiles(fileList []string) []file.File {
	results := make([]file.File, len(fileList))
	var wg sync.WaitGroup

	for i, filePath := range fileList {
		wg.Add(1)
		go func(idx int, path string) {
			defer wg.Done()
			data, err := os.ReadFile(path)
			if err != nil {
				results[idx] = file.File{
					Content: "",
					Err:     response.Error("reading %s", path),
					Ref:     path,
				}
				return
			}
			results[idx] = file.File{
				Content: string(data),
				Ref:     path,
			}
		}(i, filePath)
	}

	wg.Wait()
	return results
}

func writeFiles(files []file.File, output string, outputIsFolder bool) {
	if !outputIsFolder {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)
	errCh := make(chan error, len(files))

	for i, f := range files {
		if f.Err != nil || f.Content == "" {
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(i int, f file.File) {
			defer wg.Done()
			defer func() { <-sem }()

			defaultOutputFolder := ""
			if output == "." {
				defaultOutputFolder = "sql"
			}

			baseOut := filepath.Join(output, defaultOutputFolder)
			if i != len(files)-1 {
				baseOut = filepath.Join(baseOut, "actions")
			}

			if err := os.MkdirAll(baseOut, 0755); err != nil {
				errCh <- fmt.Errorf("creating directory %w", err)
				return
			}

			fullPath := filepath.Join(baseOut, filepath.Base(f.Ref))

			fd, err := os.Create(fullPath)
			if err != nil {
				errCh <- fmt.Errorf("creating file %w", err)
				return
			}
			defer fd.Close()

			if _, err := fd.Write([]byte(f.Content)); err != nil {
				errCh <- fmt.Errorf("writing to file %w", err)
				return
			}
		}(i, f)
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			console.Log("Error:", err)
		}
	}
}
