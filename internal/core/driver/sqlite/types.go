package sqlite

import "salvadorsru/bob/internal/models/table"

var Types = table.NewTypes(
	table.Types{
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
		Text:      "TEXT",
		Blob:      "BLOB",
		Id:        "INTEGER",
		Boolean:   "BOOLEAN",
		Date:      "DATE",
		Time:      "TIME",
		Datetime:  "DATETIME",
		Current:   "TIMESTAMP",
		Timestamp: "TIMESTAMP",
	},
)
