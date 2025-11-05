package postgre

import "salvadorsru/bob/internal/models/table"

var Types = table.NewTypes(table.Types{
	// Integers
	Int:   "INTEGER",
	Int8:  "SMALLINT",
	Int16: "SMALLINT",
	Int32: "INTEGER",
	Int64: "BIGINT",
	// Floats
	Float32: "REAL",
	Float64: "DOUBLE PRECISION",
	// Strings
	String:   "VARCHAR(255)",
	String8:  "VARCHAR(8)",
	String16: "VARCHAR(16)",
	String32: "VARCHAR(32)",
	String64: "VARCHAR(64)",
	// Other types
	Text:      "TEXT",
	Blob:      "BYTEA",
	Date:      "DATE",
	Time:      "TIME",
	Datetime:  "DATETIME",
	Id:        "SERIAL",
	Boolean:   "BOOLEAN",
	Current:   "TIMESTAMP",
	Timestamp: "TIMESTAMP",
})
