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
	motor := args[0].String()

	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		query := args[0].String()

		transpileError, tables, actions := transpiler.Transpile(motor, query)
		if transpileError != nil {
			return js.ValueOf(map[string]any{
				"error": map[string]any{
					"name":    transpileError.Name,
					"message": transpileError.Error(),
				},
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
