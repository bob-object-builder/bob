package drivers

type Type string

func UseTag(tag string) (string, []string) {
	switch Type(tag) {
	case Id:
		var tagType = string(tag)
		return tagType, []string{string(Primary), string(AutoIncrement)}
	case CreatedAt:
		var tagType = string(Date)
		return tagType, []string{string(Default), string(Now)}
	}

	return "", nil
}

const (
	Int       Type = "int"
	Int8      Type = "int8"
	Int16     Type = "int16"
	Int32     Type = "int32"
	Int64     Type = "int64"
	Uint      Type = "uint"
	Uint8     Type = "uint8"
	Uint16    Type = "uint16"
	Uint32    Type = "uint32"
	Uint64    Type = "uint64"
	Float32   Type = "float32"
	Float64   Type = "float64"
	String    Type = "string"
	String8   Type = "string8"
	String16  Type = "string16"
	String32  Type = "string32"
	String64  Type = "string64"
	Text      Type = "text"
	Blob      Type = "blob"
	Date      Type = "date"
	Time      Type = "time"
	Id        Type = "id"
	CreatedAt Type = "createdAt"
)
