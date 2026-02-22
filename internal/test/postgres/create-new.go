package postgres_test

import test_types "salvadorsru/bob/internal/test/types"

var TestCreateNew = test_types.ToTestCreateNew{
	Driver: "postgres",
	ExpectedNewUser: `
INSERT INTO users
    (name, email, age)
VALUES (
    'John Doe', 'john@mail.com', 25
);`,
	ExpectedNewBulkUser: `
INSERT INTO users
    (name, email, age)
VALUES
    ('John Doe', 'john@mail.com', 25),
    ('Jane Doe', 'jane@mail.com', 22);
`,
}
