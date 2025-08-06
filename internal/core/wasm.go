//go:build js && wasm

package core

import (
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/transpiler"
	"syscall/js"
)

func transpileWrapper(this js.Value, args []js.Value) any {
	driver := args[0].String()
	query := args[1].String()

	result := transpiler.Transpile(drivers.Motor(driver), query)
	return js.ValueOf(result)
}

func main() {
	js.Global().Set("bob", js.FuncOf(transpileWrapper))
	select {} // keep running
}
