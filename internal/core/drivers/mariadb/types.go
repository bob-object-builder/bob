package mariadb

import "salvadorsru/bob/internal/models/table"

var Types = table.NewTypes(table.Types{
	// Integers
	Int:   "INT",
	Int8:  "TINYINT",
	Int16: "SMALLINT",
	Int32: "INT",
	Int64: "BIGINT",
	// Floats
	Float32: "FLOAT",
	Float64: "DOUBLE",
	// Strings
	String:   "VARCHAR(255)",
	String8:  "VARCHAR(8)",
	String16: "VARCHAR(16)",
	String32: "VARCHAR(32)",
	String64: "VARCHAR(64)",
	// Other types
	Text:      "TEXT",
	Blob:      "BLOB",
	Date:      "DATE",
	Time:      "TIME",
	Datetime:  "DATETIME",
	Id:        "INT",
	Boolean:   "BOOLEAN",
	Current:   "TIMESTAMP",
	Timestamp: "TIMESTAMP",
})
