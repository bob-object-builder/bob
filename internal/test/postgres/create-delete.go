package postgres_test

import test_types "salvadorsru/bob/internal/test/types"

var TestCreateDelete = test_types.ToTestCreateDelete{
	Driver: "postgres",
	ExpectedDeleteUser: `
		DELETE FROM users
		WHERE
		    id = ?;
	`,
}
