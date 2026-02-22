package sqlite_test

import test_types "salvadorsru/bob/internal/test/types"

var TestCreateGet = test_types.ToTestCreateGet{
	Driver: "sqlite",
	ExpectedGetUsersWithEmailAndProfile: `
SELECT
    id,
    name,
    email
FROM users
LEFT JOIN profiles ON users.profiles_id = profiles.id
WHERE
    users.email = 'test@test.com'
    AND profiles.avatar = 'default.png';
	`,
	ExpectedGetUsersWithProfileAvatar: `
SELECT
    name,
    profiles.avatar AS profiles_avatar
FROM users
LEFT JOIN profiles ON users.profiles_id = profiles.id;
	`,
	ExpectedGetUsersComplex: `
SELECT
    id,
    name,
    email,
    family_name AS last_name,
    CONCAT('Hello ', users.name) AS meet,
    (
    SELECT
        name
    FROM profiles
    WHERE
        profiles.id = users.profiles_id
        AND profiles.avatar != 'default1.png'
    ) AS profile
FROM users
LEFT JOIN profiles ON users.profiles_id = profiles.id
LEFT JOIN post ON users.post_uuid = post.uuid
WHERE
    users.email = 'test@test.com'
    AND users.name = 'John'
    AND profiles.avatar != 'default2.png'
    AND post.upvotes > 100;
	`,
	ExpectedGetUsersGroupAge: `
SELECT
    *
FROM users
GROUP BY
    age
HAVING
    users.age > 18;
	`,
	ExpectedGetPostsGroupRating: `
SELECT
    rating,
    COUNT(posts.id) AS total_posts
FROM posts
GROUP BY
    rating
HAVING
    total_posts > 10;
	`,
}
