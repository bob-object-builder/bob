package mariadb_test

import test_types "salvadorsru/bob/internal/test/types"

var TestCreateTable = test_types.ToTestCreateTable{
	Driver: "mariadb",
	ExpectedProfilesTable: `CREATE TABLE profiles (
    id INT AUTO_INCREMENT NOT NULL,
    avatar VARCHAR(255) NOT NULL,
    PRIMARY KEY (id)
);`,
	ExpectedUsersTable: `CREATE TABLE users (
    id INT AUTO_INCREMENT NOT NULL,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    profiles_id INT,
    FOREIGN KEY (profiles_id) REFERENCES profiles(id),
    PRIMARY KEY (id)
);

CREATE INDEX idx_users_email ON users(email);`,
	ExpectedPostsTable: `CREATE TABLE posts (
	title VARCHAR(255) NOT NULL,
    content VARCHAR(255) NOT NULL,
    rating INT NOT NULL,
    users_id INT NOT NULL,
    FOREIGN KEY (users_id) REFERENCES users(id)
);`,
}
