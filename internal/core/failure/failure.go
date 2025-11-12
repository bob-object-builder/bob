package failure

import (
	"errors"
	"fmt"
)

type Failure struct {
	Name string
	fail error
}

func (f Failure) Error() string {
	return f.fail.Error()
}

const (
	IdLimitValueRequired        = "LimitValueRequired"
	IdLimitValueMustBeInteger   = "LimitValueMustBeInteger"
	IdOffsetValueMustBeInteger  = "OffsetValueMustBeInteger"
	IdMalformedQuery            = "MalformedQuery"
	IdConditionValidation       = "ConditionValidation"
	IdDeleteCondition           = "DeleteCondition"
	IdUndefinedReferenceTable   = "UndefinedReferenceTable"
	IdUndefinedReferencedColumn = "UndefinedReferencedColumn"
	IdUndefinedTypeForColumn    = "UndefinedTypeForColumn"
	IdInvalidSelectedColumn     = "InvalidSelectedColumn"
	IdColumnNotReceivingValue   = "ColumnNotReceivingValue"
	IdUndefinedToken            = "UndefinedToken"
	IdInvalidTypeForColumn      = "InvalidTypeForColumn"
	IdInvalidProperty           = "InvalidProperty"
	IdUnknownDriver             = "UnknownDriver"
	IdMalformedCondition        = "MalformedCondition"
	IdJsonParse                 = "JsonParse"
	IdMalformedArgs             = "MalformedArgs"
	IdCollectFiles              = "CollectFiles"
	IdInvalidInput              = "InvalidInput"
	IdIO
)

var (
	LimitValueRequired       = &Failure{Name: IdLimitValueRequired, fail: errors.New("limit value is required")}
	LimitValueMustBeInteger  = &Failure{Name: IdLimitValueMustBeInteger, fail: errors.New("limit value must be an integer")}
	OffsetValueMustBeInteger = &Failure{Name: IdOffsetValueMustBeInteger, fail: errors.New("offset value must be an integer")}
	JsonParse                = &Failure{Name: IdJsonParse, fail: errors.New("error during json parsing")}
	MalformedArgs            = &Failure{Name: IdMalformedArgs, fail: errors.New("error on arguments provided")}
	CollectFiles             = &Failure{Name: IdCollectFiles, fail: errors.New("error on collect files")}
	InvalidInput             = &Failure{Name: IdInvalidInput, fail: errors.New("invalid input")}
	IO                       = &Failure{Name: IdIO, fail: errors.New("io error")}
)

func MalformedQuery(token string) *Failure {
	return &Failure{
		Name: IdMalformedQuery,
		fail: fmt.Errorf("malformed query %s", token),
	}
}

func ConditionValidation(table, target string) *Failure {
	return &Failure{
		Name: IdConditionValidation,
		fail: fmt.Errorf("validation failed for table '%s' target '%s' condition cannot be empty", table, target),
	}
}

func DeleteCondition(table string) *Failure {
	return &Failure{
		Name: IdDeleteCondition,
		fail: fmt.Errorf("you must add conditions or use '*' to delete all in '%s'", table),
	}
}

func UndefinedReferenceTable(table string) *Failure {
	return &Failure{
		Name: IdUndefinedReferenceTable,
		fail: fmt.Errorf("undefined reference table '%s'", table),
	}
}

func UndefinedReferencedColumn(column, table string) *Failure {
	return &Failure{
		Name: IdUndefinedReferencedColumn,
		fail: fmt.Errorf("undefined referenced column: '%s' in table '%s'", column, table),
	}
}

func UndefinedTypeForColumn(column string) *Failure {
	return &Failure{
		Name: IdUndefinedTypeForColumn,
		fail: fmt.Errorf("undefined type for column: %s", column),
	}
}

func InvalidSelectedColumn(column string) *Failure {
	return &Failure{
		Name: IdInvalidSelectedColumn,
		fail: fmt.Errorf("invalid selected column '%s'", column),
	}
}

func ColumnNotReceivingValue(column, position string) *Failure {
	return &Failure{
		Name: IdColumnNotReceivingValue,
		fail: fmt.Errorf("column '%s' is not receiving a value at [%s]", column, position),
	}
}

func UndefinedToken(token string) *Failure {
	return &Failure{
		Name: IdUndefinedToken,
		fail: fmt.Errorf("undefined token '%s'", token),
	}
}

func InvalidTypeForColumn(column, typ string) *Failure {
	return &Failure{
		Name: IdInvalidTypeForColumn,
		fail: fmt.Errorf("invalid type %q for column '%s'", typ, column),
	}
}

func InvalidProperty(property string) *Failure {
	return &Failure{
		Name: IdInvalidProperty,
		fail: fmt.Errorf("invalid property '%s'", property),
	}
}

func UnknownDriver(driver string) *Failure {
	return &Failure{
		Name: IdUnknownDriver,
		fail: fmt.Errorf("unknown driver: %s", driver),
	}
}

func MalformedCondition(condition string) *Failure {
	error := fmt.Errorf("malformed condition")
	if condition != "" {
		error = fmt.Errorf("malformed condition %s", condition)
	}

	return &Failure{
		Name: IdMalformedCondition,
		fail: error,
	}
}
