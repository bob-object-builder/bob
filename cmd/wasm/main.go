//go:build js && wasm

package main

import (
	"salvadorsru/bob/internal/core/transpiler"
	"syscall/js"
)

func sliceToJSArray(slice []string) js.Value {
	jsArray := js.Global().Get("Array").New()
	for _, s := range slice {
		jsArray.Call("push", s)
	}
	return jsArray
}

func bob(this js.Value, args []js.Value) any {
	driverString := args[0].String()
	driverError, driver := transpiler.GetDriver(driverString)
	if driverError != nil {
		return js.ValueOf(map[string]any{
			"error": driverError.Error(),
			"value": nil,
		})
	}

	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		query := args[0].String()

		transpileError, tables, actions := transpiler.Transpile(driver, query)
		if transpileError != nil {
			return js.ValueOf(map[string]any{
				"error": transpileError.Error(),
				"value": nil,
			})
		}

		return js.ValueOf(map[string]any{
			"error": nil,
			"value": map[string]any{
				"tables":  sliceToJSArray(tables.Get()),
				"actions": sliceToJSArray(actions.Get()),
			},
		})
	})

	return fn
}

func main() {
	js.Global().Set("bob", js.FuncOf(bob))
	select {}
}
