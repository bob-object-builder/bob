package mariadb_test

import test_types "salvadorsru/bob/internal/test/types"

var TestCreateDelete = test_types.ToTestCreateDelete{
	Driver: "mariadb",
	ExpectedDeleteUser: `
		DELETE FROM users
		WHERE
		    id = ?;
	`,
}
