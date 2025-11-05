package table

import (
	"fmt"
)

type Type string

type TypesMap map[Type]Type

const IdType Type = "id"
const CurrentType Type = "current"

type Types struct {
	Int, Int8, Int16, Int32, Int64                Type
	Float32, Float64                              Type
	String, String8, String16, String32, String64 Type
	Date, Time, Datetime, Timestamp, Current      Type
	Text, Blob, Id, Boolean                       Type
}

func NewTypes(t Types) TypesMap {
	return TypesMap{
		"int":       t.Int,
		"int8":      t.Int8,
		"int16":     t.Int16,
		"int32":     t.Int32,
		"int64":     t.Int64,
		"float32":   t.Float32,
		"float64":   t.Float64,
		"string":    t.String,
		"string8":   t.String8,
		"string16":  t.String16,
		"string32":  t.String32,
		"string64":  t.String64,
		"text":      t.Text,
		"blob":      t.Blob,
		"id":        t.Id,
		"boolean":   t.Boolean,
		"time":      t.Time,
		"date":      t.Date,
		"datetime":  t.Datetime,
		"current":   t.Current,
		"timestamp": t.Timestamp,
	}
}

var typeMap = TypesMap{
	"int":       "int",
	"int8":      "int8",
	"int16":     "int16",
	"int32":     "int32",
	"int64":     "int64",
	"float32":   "float32",
	"float64":   "float64",
	"string":    "string",
	"string8":   "string8",
	"string16":  "string16",
	"string32":  "string32",
	"string64":  "string64",
	"text":      "text",
	"blob":      "blob",
	"id":        "id",
	"boolean":   "boolean",
	"time":      "time",
	"date":      "date",
	"datetime":  "datetime",
	"current":   "current",
	"timestamp": "timestamp",
}

func (t TypesMap) GetType(token Type) (Type, error) {
	if typ, ok := t[token]; ok {
		return typ, nil
	}
	return "", fmt.Errorf("invalid type %q", token)
}

func IsType(token string) bool {
	if typ, ok := typeMap[Type(token)]; ok {
		return typ != ""
	}
	return false
}
