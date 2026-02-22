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

func TestCreateNew(t *testing.T) {
	tests := value.NewArray[test_types.ToTestCreateNew](
		sqlite_test.TestCreateNew,
		postgres_test.TestCreateNew,
		mariadb_test.TestCreateNew,
	)

	for _, test := range *tests {
		transpileError, _, actions := transpiler.Transpile(test.Driver, `
			new Users {
			    name    "John Doe"
			    email   "john@mail.com"
			    age     25
			}

			new Users name email age {
			    "John Doe" "john@mail.com" 25
			    "Jane Doe" "jane@mail.com" 22
			}
		`)

		if transpileError != nil {
			t.Fatalf("Driver %s: Transpile error: %s", test.Driver, transpileError.Message)
		}

		AssertStringsEqual(t, test.Driver, actions[0], test.ExpectedNewUser)
		AssertStringsEqual(t, test.Driver, actions[1], test.ExpectedNewBulkUser)
	}
}
