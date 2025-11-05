package console

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func Log(v ...any) {
	format := strings.TrimSpace(strings.Repeat("%v ", len(v)))
	fmt.Printf(format+"\n", v...)
}

func Success(v ...any) {
	// Log("success: " + response.Success(v...))
	Log("success")
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
