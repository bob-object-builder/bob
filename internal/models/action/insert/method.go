package insert

import (
	"fmt"
	"salvadorsru/bob/internal/core/console"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func (new New) ToQuery(driver drivers.Driver) string {
	query := "INSERT INTO `%s`\n%s\nVALUES\n%s"

	values := []string{}

	for _, v := range new.Values {
		values = append(values, "("+strings.Join(v, ", ")+")")
	}

	console.Log()

	return fmt.Sprintf(query, new.Table, utils.Indent("("+strings.Join(new.Fields, ", ")+")"), utils.IndentLines(strings.Join(values, ",\n")))
}
