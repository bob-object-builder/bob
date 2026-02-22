package test

import (
	"salvadorsru/bob/internal/core/transpiler"
	"salvadorsru/bob/internal/lib/value"
	mariadb_test "salvadorsru/bob/internal/test/mariadb"
	postgres_test "salvadorsru/bob/internal/test/postgres"
	sqlite_test "salvadorsru/bob/internal/test/sqlite"
	test_types "salvadorsru/bob/internal/test/types"
	"testing"
)

func TestCreateDelete(t *testing.T) {
	tests := value.NewArray[test_types.ToTestCreateDelete](
		sqlite_test.TestCreateDelete,
		postgres_test.TestCreateDelete,
		mariadb_test.TestCreateDelete,
	)

	for _, test := range *tests {
		transpileError, _, actions := transpiler.Transpile(test.Driver, `
			delete Users {
				if id = ?
			}
		`)

		if transpileError != nil {
			t.Fatalf("Driver %s: Transpile error: %s", test.Driver, transpileError.Message)
		}

		AssertStringsEqual(t, test.Driver, actions[0], test.ExpectedDeleteUser)
	}
}
