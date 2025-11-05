package transpiler

import (
	"fmt"
	"strings"

	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value/array"
	"salvadorsru/bob/internal/models/condition"
	"salvadorsru/bob/internal/models/get"
)

const EveryField = "*"
const SpreadEveryField = "..."

func (t Transpiler) TranspileGet(g get.Get) string {
	queryTemplate := "SELECT\n%s\nFROM %s%s%s%s%s;"

	fields := t.transpileFields(g)
	joins, conditions, groups, having := t.collectJoinsAndConditionsAndHaving(g)

	selectedString := formatter.IndentLines(strings.Join(*fields, ",\n"))
	joinString := ""
	if len(*joins) > 0 {
		joinString = "\n" + formatter.IndentLines(strings.Join(*joins, "\n"))
	}

	conditionString := ""
	if len(*conditions) > 0 {
		conditionString = "\nWHERE\n" + t.TranspileConditions(*conditions)
	}

	groupString := ""
	if len(*groups) > 0 {
		groupString = "\nGROUP BY\n" + formatter.IndentLines(strings.Join(*groups, ",\n"))
	}

	havingString := ""
	if len(*having) > 0 {
		havingString = "\nHAVING\n" + t.TranspileConditions(*having)
	}

	return fmt.Sprintf(queryTemplate, selectedString, g.Target, joinString, conditionString, groupString, havingString)
}

func (t Transpiler) transpileFields(g get.Get) *array.Array[string] {
	fields := array.New[string]()

	for field := range g.Selected.Range() {
		if field.Key == field.Value || field.Value == EveryField || field.Value == SpreadEveryField {

			if field.Key == SpreadEveryField || field.Value == SpreadEveryField {
				field.Value = EveryField
			}

			fields.Push(field.Value)
		} else {
			fields.Push(fmt.Sprintf("%s AS %s", field.Value, field.Key))
		}
	}

	if fields.Length() == 0 {
		fields.Push(EveryField)
	}

	for _, sub := range g.Subqueries {
		sql := t.TranspileGet(sub)
		fields.Push(
			fmt.Sprintf("(\n%s\n) AS %s", formatter.IndentLines(sql), sub.Alias),
		)
	}

	for _, join := range g.Joins {
		for field := range join.Selected.Range() {
			alias := field.Key
			if field.Value == field.Key {
				alias = fmt.Sprintf("%s_%s", join.Target, field.Value)
			}
			fields.Push(fmt.Sprintf("%s AS %s", field.Value, alias))
		}
	}

	return fields
}

func (t Transpiler) collectJoinsAndConditionsAndHaving(g get.Get) (
	*array.Array[string],
	*array.Array[condition.Condition],
	*array.Array[string],
	*array.Array[condition.Condition],
) {
	joins := array.New[string]()
	conditions := array.New[condition.Condition](g.Conditions...)
	groups := array.New[string](g.Groups...)
	having := array.New[condition.Condition](g.Having...)

	for _, join := range g.Joins {
		joins.Push(t.TranspileLeftJoin(g.Target, join))
		conditions.Push(join.Conditions...)
		having.Push(join.Having...)
		groups.Push(join.Groups...)
	}

	return joins, conditions, groups, having
}
