package file

import (
	"fmt"
	"os"
	"path/filepath"
	"salvadorsru/bob/internal/lib/console"
	"salvadorsru/bob/internal/lib/response"
	"sync"
)

type File struct {
	Content string
	Err     error
	Ref     string
}

func (f *File) ToString() (string, error) {
	if f.Err != nil {
		return "", response.Error("reading %s", f.Ref)
	}
	return f.Content, nil
}

func FilesToString(f []File) (string, error) {
	var content string
	for _, file := range f {
		if file.Err != nil {
			return "", response.Error("reading %s", file.Ref)
		}
		content += file.Content
	}
	return content, nil
}

func WriteFiles(files []File, output string, outputIsFolder bool) {
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

		go func(i int, f File) {
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

func FindBobFiles(directory string) ([]string, error) {
	var files []string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".bob" {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func ReadFiles(fileList []string) []File {
	results := make([]File, len(fileList))
	var wg sync.WaitGroup

	for i, filePath := range fileList {
		wg.Add(1)
		go func(idx int, path string) {
			defer wg.Done()
			data, err := os.ReadFile(path)
			if err != nil {
				results[idx] = File{
					Content: "",
					Err:     response.Error("reading %s", path),
					Ref:     path,
				}
				return
			}
			results[idx] = File{
				Content: string(data),
				Ref:     path,
			}
		}(i, filePath)
	}

	wg.Wait()
	return results
}
