package mariadb

import (
	"salvadorsru/bob/internal/models/function"
)

var Callers = function.NewCallers(function.Callers{
	Length: "CHAR_LENGTH",
})
