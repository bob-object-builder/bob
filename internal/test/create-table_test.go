package test

import (
	"salvadorsru/bob/internal/core/transpiler"
	"salvadorsru/bob/internal/lib/value"
	mariadb_test "salvadorsru/bob/internal/test/mariadb"
	postgres_test "salvadorsru/bob/internal/test/postgres"
	sqlite_test "salvadorsru/bob/internal/test/sqlite"
	"testing"
)

func TestCreateTable(t *testing.T) {
	tests := value.NewArray(
		sqlite_test.TestCreateTable,
		postgres_test.TestCreateTable,
		mariadb_test.TestCreateTable,
	)

	for _, test := range *tests {
		transpileError, tables, _ := transpiler.Transpile(test.Driver, `
			table Profiles {
				id id
				avatar string
			}

			table Users {
				id id
				name string16
				email string unique index
				created_at current
				Profiles id optional
				age int int16

				if age > 18
			}

			table Posts {
				title string
				content string
				rating int
				Users
			}
		`)

		if transpileError != nil {
			t.Fatalf("Driver %s: Transpile error: %s", test.Driver, transpileError.Message)
		}

		AssertStringsEqual(t, test.Driver, tables[0], test.ExpectedProfilesTable)
		AssertStringsEqual(t, test.Driver, tables[1], test.ExpectedUsersTable)
		AssertStringsEqual(t, test.Driver, tables[2], test.ExpectedPostsTable)
	}
}
