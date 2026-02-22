package sqlite_test

import test_types "salvadorsru/bob/internal/test/types"

var TestCreateDelete = test_types.ToTestCreateDelete{
	Driver: "sqlite",
	ExpectedDeleteUser: `
		DELETE FROM users
		WHERE
		    id = ?;
	`,
}
