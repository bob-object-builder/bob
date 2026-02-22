package sqlite_test

import test_types "salvadorsru/bob/internal/test/types"

var TestCreateTable = test_types.ToTestCreateTable{
	Driver: "sqlite",
	ExpectedProfilesTable: `
	CREATE TABLE profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			avatar TEXT NOT NULL
	);`,
	ExpectedUsersTable: `CREATE TABLE users (
	  id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
	  name TEXT NOT NULL,
	  email TEXT NOT NULL,
	  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	  profiles_id INTEGER,
	  FOREIGN KEY (profiles_id) REFERENCES profiles(id)
	);

	CREATE INDEX idx_users_email ON users(email);`,
	ExpectedPostsTable: `
	CREATE TABLE posts (
	  title TEXT NOT NULL,
	  content TEXT NOT NULL,
	  rating INTEGER NOT NULL,
	  users_id INTEGER NOT NULL,
	  FOREIGN KEY (users_id) REFERENCES users(id)
	);`,
}
