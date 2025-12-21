package driver

type Type string

const (
	Int       Type = "int"
	Int8      Type = "int8"
	Int16     Type = "int16"
	Int32     Type = "int32"
	Int64     Type = "int64"
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
	Datetime  Type = "datetime"
	Id        Type = "id"
	Boolean   Type = "boolean"
	Current   Type = "current"
	Timestamp Type = "timestamp"
)

type Types map[Type]string

func (t *Types) Get(token string) string {
	if val, ok := (*t)[Type(token)]; ok {
		return val
	}

	return ""
}

type Map struct {
	Int, Int8, Int16, Int32, Int64                string
	Float32, Float64                              string
	String, String8, String16, String32, String64 string
	Text, Blob                                    string
	Date, Time, Datetime                          string
	Id, Boolean, Current, Timestamp               string
}

func IsType(t string) bool {
	switch Type(t) {
	case Int, Int8, Int16, Int32, Int64, Float32, Float64, String, String8, String16, String32, String64, Text, Blob, Date, Time, Datetime, Id, Boolean, Current, Timestamp:
		return true
	default:
		return false
	}
}

func NewTypes(m Map) Types {
	return Types{
		Int:       m.Int,
		Int8:      m.Int8,
		Int16:     m.Int16,
		Int32:     m.Int32,
		Int64:     m.Int64,
		Float32:   m.Float32,
		Float64:   m.Float64,
		String:    m.String,
		String8:   m.String8,
		String16:  m.String16,
		String32:  m.String32,
		String64:  m.String64,
		Text:      m.Text,
		Blob:      m.Blob,
		Date:      m.Date,
		Time:      m.Time,
		Datetime:  m.Datetime,
		Id:        m.Id,
		Boolean:   m.Boolean,
		Current:   m.Current,
		Timestamp: m.Timestamp,
	}
}
