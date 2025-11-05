//go:build js && wasm

package main

import (
	"fmt"
	"salvadorsru/bob/internal/core/transpiler"
	"syscall/js"
)

func joinTemplateLiteral(args []js.Value) string {
	if len(args) == 0 {
		return ""
	}

	stringsArr := args[0]
	if stringsArr.Type() != js.TypeObject || !stringsArr.InstanceOf(js.Global().Get("Array")) {
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
	driverString := args[0].String()
	driverError, driver := transpiler.GetDriver(driverString)
	if driverError != nil {
		return js.ValueOf(map[string]any{
			"error": driverError.Error(),
			"value": nil,
		})
	}

	fn := js.FuncOf(func(this js.Value, args []js.Value) any {
		query := joinTemplateLiteral(args)

		transpileError, tables, actions := transpiler.Transpile(driver, query)
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
