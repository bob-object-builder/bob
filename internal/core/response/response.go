package response

import (
	"fmt"
	"strings"
)

func Success(v ...any) string {
	toPrint := strings.TrimSpace(strings.Repeat("%s ", len(v)))
	return fmt.Sprintf(toPrint, v...)
}

func Error(v ...any) error {
	toPrint := strings.TrimSpace(strings.Repeat("%s ", len(v)))
	return fmt.Errorf(toPrint, v...)
}
