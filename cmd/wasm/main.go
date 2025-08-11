//go:build js && wasm

package main

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/transpiler"
	"syscall/js"
)

func join_template_literal(args []js.Value) string {
	if len(args) == 0 {
		return ""
	}

	stringsArr := args[0]
	if stringsArr.Type() != js.TypeObject || !stringsArr.InstanceOf(js.Global().Get("Array")) {
		// Not an array, assuming plain string
		return args[0].String()
	}

	length := stringsArr.Length()
	result := ""

	for i := 0; i < int(length); i++ {
		result += stringsArr.Index(i).String()
		if i+1 < len(args) {
			result += args[i+1].String()
		}
	}

	return result
}

func bob(this js.Value, args []js.Value) any {
	driver := args[0].String()

	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		query := join_template_literal(args)

		transpileError, tables, actions := transpiler.Transpile(drivers.Motor(driver), query)
		if transpileError != nil {
			return js.ValueOf(map[string]any{
				"error": transpileError.Error(),
				"value": nil,
			})
		}

		result := fmt.Sprintf("%s\n\n%s\n", tables, actions)
		return js.ValueOf(map[string]any{
			"error": nil,
			"value": result,
		})
	})

	return fn
}

func main() {
	js.Global().Set("bob", js.FuncOf(bob))
	select {} // keep running
}
