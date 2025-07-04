package table

import (
	"fmt"
	"salvadorsru/bob/internal/core/drivers"
	"salvadorsru/bob/internal/core/utils"
	"strings"
)

func makeIndexSentence(indexableList ...Indexable) string {
	sentences := []string{}

	for _, indexable := range indexableList {
		sentence := indexable.ToQuery()
		sentences = append(sentences, sentence)
	}

	return strings.Join(sentences, ";\n")
}

func makeColumn(driver drivers.Driver, t *Table, columnName string, hasOnlyOnePrimeryKey bool) string {
	column := t.Columns.Get(columnName)

	columnType := driver.GetType(column.Type)

	columnSentence := fmt.Sprintf("%s %s", columnName, columnType)

	if hasOnlyOnePrimeryKey && column.IsPrimaryKey {
		columnSentence += " PRIMARY KEY"
	}

	if column.IsAutoIncrement {
		attribute := driver.GetAttribute(string(drivers.AutoIncrement))
		if attribute != "" {
			columnSentence += " " + attribute
		}
	}

	if column.Default != nil {

		literal := driver.GetLiteral(*column.Default)
		if literal != "" {
			columnSentence += fmt.Sprintf(" DEFAULT %v", literal)
		} else {
			columnSentence += fmt.Sprintf(" DEFAULT %v", *column.Default)
		}
	}

	columnSentence = utils.Indent(columnSentence)

	return columnSentence
}

func (t Table) ToQuery(driver drivers.Driver) string {
	query := "CREATE TABLE %s (\n%s\n)"
	columns := []string{}
	hasOnlyOnePrimeryKey := len(t.PrimaryKeys) == 1
	tableName := t.Name

	for _, columnName := range t.Columns.Order {
		columnSentence := makeColumn(driver, &t, columnName, hasOnlyOnePrimeryKey)
		columns = append(columns, columnSentence)
	}

	for _, reference := range t.References.Data {
		reference := reference.toQuery()
		reference = utils.Indent(reference)
		columns = append(columns, reference)
	}

	if !hasOnlyOnePrimeryKey && len(t.PrimaryKeys) > 1 {
		primaryKeysSentence := fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(t.PrimaryKeys, ", "))
		primaryKeysSentence = utils.Indent(primaryKeysSentence)
		columns = append(columns, primaryKeysSentence)
	}

	if len(t.Uniques) > 0 {
		uniqueSentence := fmt.Sprintf("UNIQUE (%s)", strings.Join(t.Uniques, ", "))
		uniqueSentence = utils.Indent(uniqueSentence)
		columns = append(columns, uniqueSentence)
	}

	indexSentences := makeIndexSentence(t.Indexes...)

	columnsSentence := strings.Join(columns, ",\n")

	query = fmt.Sprintf(query, tableName, columnsSentence)

	if indexSentences != "" {
		query += ";\n\n" + indexSentences
	}

	return query
}
