package console

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"salvadorsru/bob/internal/core/response"
	"strings"
)

func Log(v ...any) {
	toPrint := strings.TrimSpace(strings.Repeat("%s ", len(v)))
	fmt.Printf(toPrint+"\n", v...)
}

func Success(v ...any) {
	Log("success: " + response.Success(v...))
}

func Panic(v ...any) {
	toPrint := strings.TrimSpace(strings.Repeat("%s ", len(v)))
	fmt.Printf("\033[31m"+"error: "+toPrint+"\033[0m\n", v...)
	os.Exit(1)
}

func Clear() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
