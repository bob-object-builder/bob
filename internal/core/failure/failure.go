package failure

import "fmt"

type Failure struct {
	Name    string
	Message string
}

/*
|--------------------------------------------------------------------------
| Failure IDs
|--------------------------------------------------------------------------
*/

const (
	IdUndefinedTypeForColumn               = "UndefinedTypeForColumn"
	IdConditionValidation                  = "ConditionValidation"
	IdUndefinedConditionComparator         = "UndefinedConditionComparator"
	IdMalformedCondition                   = "MalformedCondition"
	IdUndefinedFunction                    = "UndefinedFunction"
	IdInvalidMotor                         = "InvalidMotor"
	IdUndefinedKeyword                     = "UndefinedKeyword"
	IdTypeNotFound                         = "TypeNotFound"
	IdNotEnoughValues                      = "NotEnoughValues"
	IdUnknownDriver                        = "UnknownDriver"
	IdUndefinedValueOnSetter               = "UndefinedValueOnSetter"
	IdUndefinedNullsDefinition             = "UndefinedNullsDefinition"
	IdSubqueryMustBeSingleColumn           = "SubqueryMustBeSingleColumn"
	IdCannotUseExpressionOrStringSelection = "CannotUseExpressionOrStringAsSelection"
	IdTypeDoesNotImplementMerge            = "TypeDoesNotImplementMerge"
	IdUnclosedExpression                   = "UnclosedExpression"
	IdUnclosedStringLiteral                = "UnclosedStringLiteral"
	IdInvalidExpressionBeforeParenthesis   = "InvalidExpressionBeforeParenthesis"
	IdUnexpectedClosingParenthesis         = "UnexpectedClosingParenthesis"
	IdExpressionProcessingEmptyResult      = "ExpressionProcessingEmptyResult"
	IdTypeDoesNotImplementTableChild       = "TypeDoesNotImplementTableChild"
	IdMissingAliasForExpression            = "MissingAliasForExpression"
	IdSelectionValueIsAlreadyDefined       = "SelectionValueIsAlreadyDefined"
	IdInvalidTableColumn                   = "InvalidTableColumn"
	IdInvalidTableReference                = "InvalidTableReference"
	IdTableNameIsEmpty                     = "TableNameIsEmpty"
	IdInvalidInput                         = "InvalidInput"
	IdCollectFiles                         = "CollectFiles"
	IdIO                                   = "IO"
	IdMalformedArgs                        = "MalformedArgs"
	IdJsonParse                            = "JsonParse"
	IdUndefinedOrderTarget                 = "UndefinedOrderTarget"
)

/*
|--------------------------------------------------------------------------
| Static Failures
|--------------------------------------------------------------------------
*/

var (
	UndefinedNullsDefinition = &Failure{
		Name:    IdUndefinedNullsDefinition,
		Message: "undefined nulls definition",
	}
	SubqueryMustBeSingleColumn = &Failure{
		Name:    IdSubqueryMustBeSingleColumn,
		Message: "subquery must select a single column",
	}
	CannotUseExpressionOrStringAsSelection = &Failure{
		Name:    IdCannotUseExpressionOrStringSelection,
		Message: "cannot use expression or string as selection",
	}
	TypeDoesNotImplementMerge = &Failure{
		Name:    IdTypeDoesNotImplementMerge,
		Message: "type does not implement Merge method",
	}
	UnclosedExpression = &Failure{
		Name:    IdUnclosedExpression,
		Message: "unclosed expression because of missing ')'",
	}
	UnclosedStringLiteral = &Failure{
		Name:    IdUnclosedStringLiteral,
		Message: "unclosed string literal",
	}
	InvalidExpressionBeforeParenthesis = &Failure{
		Name:    IdInvalidExpressionBeforeParenthesis,
		Message: "invalid expression before '('",
	}
	UnexpectedClosingParenthesis = &Failure{
		Name:    IdUnexpectedClosingParenthesis,
		Message: "unexpected ')' without matching '('",
	}
	ExpressionProcessingEmptyResult = &Failure{
		Name:    IdExpressionProcessingEmptyResult,
		Message: "expression processing error empty result",
	}
	TypeDoesNotImplementTableChild = &Failure{
		Name:    IdTypeDoesNotImplementTableChild,
		Message: "type does not implement TableChild",
	}
	MissingAliasForExpression = &Failure{
		Name:    IdMissingAliasForExpression,
		Message: "missing alias for expression",
	}
	SelectionValueIsAlreadyDefined = &Failure{
		Name:    IdSelectionValueIsAlreadyDefined,
		Message: "selection value is already defined",
	}
	InvalidTableColumn = &Failure{
		Name:    IdInvalidTableColumn,
		Message: "invalid table column",
	}
	InvalidTableReference = &Failure{
		Name:    IdInvalidTableReference,
		Message: "invalid table reference",
	}
	TableNameIsEmpty = &Failure{
		Name:    IdTableNameIsEmpty,
		Message: "table name is empty",
	}
	InvalidInput = &Failure{
		Name:    IdInvalidInput,
		Message: "invalid input",
	}
	CollectFiles = &Failure{
		Name:    IdCollectFiles,
		Message: "error on collect files",
	}
	IO = &Failure{
		Name:    IdIO,
		Message: "io error",
	}
	MalformedArgs = &Failure{
		Name:    IdMalformedArgs,
		Message: "error on arguments provided",
	}
	JsonParse = &Failure{
		Name:    IdJsonParse,
		Message: "error during json parsing",
	}
	UndefinedOrderTarget = &Failure{
		Name:    IdUndefinedOrderTarget,
		Message: "undefined order target",
	}
)

/*
|--------------------------------------------------------------------------
| Dynamic Failures
|--------------------------------------------------------------------------
*/

func UndefinedTypeForColumn(tp string, column string) *Failure {
	return &Failure{
		Name:    IdUndefinedTypeForColumn,
		Message: fmt.Sprintf("undefined type '%s' for column '%s'", tp, column),
	}
}

func UndefinedConditionComparator(key string) *Failure {
	return &Failure{
		Name:    IdUndefinedConditionComparator,
		Message: fmt.Sprintf("undefined condition comparator '%s'", key),
	}
}

func ConditionValidation(table, target string) *Failure {
	return &Failure{
		Name: IdConditionValidation,
		Message: fmt.Sprintf(
			"validation failed for condition in '%s' with target '%s' condition cannot be empty",
			table,
			target,
		),
	}
}

func MalformedCondition(condition string) *Failure {
	message := "malformed condition"
	if condition != "" {
		message = fmt.Sprintf("malformed condition %s", condition)
	}

	return &Failure{
		Name:    IdMalformedCondition,
		Message: message,
	}
}

func UndefinedFunction(function string) *Failure {
	return &Failure{
		Name:    IdUndefinedFunction,
		Message: fmt.Sprintf("undefined function '%s'", function),
	}
}

func InvalidMotor(motor string) *Failure {
	return &Failure{
		Name:    IdInvalidMotor,
		Message: fmt.Sprintf("invalid motor '%s'", motor),
	}
}

func UndefinedKeyword(keyType string, keyword string) *Failure {
	return &Failure{
		Name:    IdUndefinedKeyword,
		Message: fmt.Sprintf("undefined '%s' keyword '%s'", keyType, keyword),
	}
}

func TypeNotFound(typeName string) *Failure {
	return &Failure{
		Name:    IdTypeNotFound,
		Message: fmt.Sprintf("type not found '%s'", typeName),
	}
}

func NotEnoughValues(at string) *Failure {
	return &Failure{
		Name:    IdNotEnoughValues,
		Message: fmt.Sprintf("not enough values for '%s' insertion", at),
	}
}

func UnknownDriver(driver string) *Failure {
	message := "empty driver"
	if driver != "" {
		message = fmt.Sprintf("unknown driver '%s'", driver)
	}

	return &Failure{
		Name:    IdUnknownDriver,
		Message: message,
	}
}

func UndefinedValueOnSetter(setter string, key string) *Failure {
	return &Failure{
		Name:    IdUndefinedValueOnSetter,
		Message: fmt.Sprintf("undefined value on setter '%s' on key '%s'", setter, key),
	}
}
