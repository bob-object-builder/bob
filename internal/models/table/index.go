package table

import (
	"fmt"
	"strings"
)

type Indexable struct {
	tableName string
	list      []string
}

func (i *Indexable) ToQuery() string {
	if len(i.list) > 0 {
		columnsNameSentence := strings.Join(i.list, "_")
		columnsNameAsParameterSentence := strings.Join(i.list, ", ")
		indexSentence := fmt.Sprintf("CREATE INDEX idx_%s_%s ON %s(%s)", i.tableName, columnsNameSentence, i.tableName, columnsNameAsParameterSentence)
		return indexSentence
	}

	return ""
}

func NewIndexable(tableName string, columns ...string) Indexable {
	return Indexable{tableName, columns}
}
