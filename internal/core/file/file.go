package file

import (
	"salvadorsru/bob/internal/core/response"
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
