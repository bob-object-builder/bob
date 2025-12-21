package transpiler

import (
	"fmt"
	"salvadorsru/bob/internal/core/failure"
	"salvadorsru/bob/internal/lib/formatter"
	"salvadorsru/bob/internal/lib/value"
	"salvadorsru/bob/internal/model/get"
	"strings"
)

func TranspileSelection(s *get.Selection, alias string) string {
	if s.From == "" || s.IsExpression {
		return fmt.Sprintf("%s AS %s", s.Selected.Value, alias)
	}

	return fmt.Sprintf("%s.%s AS %s_%s", s.From, s.Selected.Value, s.From, s.Selected.Value)
}

func TranspileGet(g *get.Get, isSubquery bool) (*failure.Failure, string) {
	selected := value.NewArray[string]()
	orders := value.NewArray[string]()

	joinsError, joinsQueries, joinsSelections := TranspileJoin(&g.Joins, g.Target)
	if joinsError != nil {
		return joinsError, ""
	}

	for _, joinSelection := range joinsSelections {
		if joinSelection.IsExpression && joinSelection.Alias == "" {
			return failure.MissingAliasForExpression, ""
		}

		g.Children.Push(get.Child{
			Alias: joinSelection.Alias,
			Value: &get.Selection{
				IsExpression: joinSelection.IsExpression,
				From:         joinSelection.From,
				Selected: value.Token{
					Type:  value.Text,
					Value: joinSelection.Target,
				},
			},
		})
	}

	for _, col := range g.Children {
		switch v := col.Value.(type) {
		case *get.Selection:
			q := TranspileSelection(v, col.Alias)
			selected.Push(formatter.Indent(q))

		case *get.Get:
			subErr, subSQL := TranspileGet(v, true)
			if subErr != nil {
				return subErr, ""
			}
			q := fmt.Sprintf("(\n%s\n) AS %s", formatter.IndentLines(subSQL), col.Alias)
			selected.Push(formatter.IndentLines(formatter.ToReferenceCase(q)))

		default:
			selected.Push(formatter.Indent(formatter.ToReferenceCase(col.Alias)))
		}
	}

	if !isSubquery {
		hasWildcard := false
		var normalized []string

		for _, s := range *selected {
			raw := strings.TrimSpace(strings.TrimLeft(s, "\t "))
			if raw == "..." {
				raw = "*"
				s = formatter.Indent("*")
			}

			if raw == "*" {
				if hasWildcard {
					continue
				}
				hasWildcard = true
			}

			normalized = append(normalized, s)
		}

		if len(normalized) == 0 {
			normalized = append(normalized, formatter.Indent("*"))
		}

		selected = value.NewArray[string]()
		for _, s := range normalized {
			selected.Push(s)
		}
	}

	conditionsError, conditionString := TranspileCondition(&g.Conditions, false)
	if conditionsError != nil {
		return conditionsError, ""
	}

	var groupString string
	if g.HasGroup {
		groupString = TranspileGroup(&g.Group)

		havingsError, havingString := TranspileCondition(&g.Conditions, true)
		if havingsError != nil {
			return havingsError, ""
		}

		if havingString != "" {
			if !strings.HasSuffix(groupString, "\n") {
				groupString += "\n"
			}
			groupString += havingString
		}
	}

	for _, o := range g.Orders {
		dir := strings.ToUpper(o.Direction)
		var nullOrder string
		if o.NullsFirst {
			nullOrder = fmt.Sprintf("(%s.%s IS NOT NULL), %s.%s %s", g.Target, o.Target, g.Target, o.Target, dir)
		} else {
			nullOrder = fmt.Sprintf("(%s.%s IS NULL), %s.%s %s", g.Target, o.Target, g.Target, o.Target, dir)
		}
		orders.Push(formatter.Indent(nullOrder))
	}

	var ordersString string
	if orders.Length() > 0 {
		ordersString = orders.Join("\n")
		ordersString = "ORDER BY\n" + ordersString
	}

	if isSubquery && len(*selected) != 1 {
		return failure.SubqueryMustBeSingleColumn, ""
	}

	var sb strings.Builder

	sb.WriteString("SELECT\n")
	sb.WriteString(selected.Join(",\n"))
	sb.WriteString("\n")

	if g.Target != "" && !strings.Contains(joinsQueries, "\nFROM ") && !strings.HasPrefix(strings.TrimSpace(joinsQueries), "FROM ") {
		sb.WriteString("FROM ")
		sb.WriteString(g.Target)
		sb.WriteString("\n")
	}

	addClause := func(clause string) {
		if clause == "" {
			return
		}
		trimmed := strings.TrimPrefix(clause, "\n")
		if sb.Len() > 0 && !strings.HasSuffix(sb.String(), "\n") {
			sb.WriteString("\n")
		}
		sb.WriteString(trimmed)
		if !strings.HasSuffix(trimmed, "\n") {
			sb.WriteString("\n")
		}
	}

	addClause(joinsQueries)
	addClause(conditionString)

	if groupString != "" {
		addClause(groupString)
	}

	if ordersString != "" {
		addClause(ordersString + "\n")
	}

	if !isSubquery {
		finalSQL := strings.TrimRight(sb.String(), "\n ")
		finalSQL = finalSQL + ";"
		return nil, finalSQL
	}

	finalSQL := strings.TrimRight(sb.String(), "\n ")
	return nil, finalSQL
}
