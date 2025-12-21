package sqlite

import (
	"salvadorsru/bob/internal/core/drivers/driver"
)

var Types = driver.NewTypes(driver.Map{
	// Integers
	Int:   "INTEGER",
	Int8:  "INTEGER",
	Int16: "INTEGER",
	Int32: "INTEGER",
	Int64: "INTEGER",
	// Floats
	Float32: "REAL",
	Float64: "REAL",
	// Strings
	String:   "TEXT",
	String8:  "TEXT",
	String16: "TEXT",
	String32: "TEXT",
	String64: "TEXT",
	// Other types
	Text:     "TEXT",
	Blob:     "BLOB",
	Date:     "DATE",
	Time:     "TIME",
	Datetime: "DATETIME",
	Id:       "INTEGER",
	Boolean:  "INTEGER",
	Current:  "TIMESTAMP",
})
