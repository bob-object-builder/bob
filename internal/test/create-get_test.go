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

func TestCreateGet(t *testing.T) {
	tests := value.NewArray[test_types.ToTestCreateGet](
		sqlite_test.TestCreateGet,
		postgres_test.TestCreateGet,
		mariadb_test.TestCreateGet,
	)

	for _, test := range *tests {
		transpileError, _, actions := transpiler.Transpile(test.Driver, `
			get Users {
			  id
			  name
			  email
			  if email = "test@test.com"

			  -> Profiles {
			    if avatar = "default.png"
			  }
			}

			get Users {
			  name
			  -> Profiles {
			    avatar
			  }
			}

			get Users {
			  id
			  name
			  email

			  last_name: family_name

			  meet: concat("Hello ", name)

			  if email = "test@test.com"

			  profile: get Profiles {
			    name
			    if id = Users->Profiles.id
			    if avatar != "default1.png"
			  }

			  -> Profiles {
			    if avatar != "default2.png"
			  }

			  -> Post uuid {
			    if upvotes > 100
			  }

			  if name = "John"
			}

			get Users {
			  group age
			  if age > 18
			}

			get Posts {
			  rating
			  total_posts: count(id)
			  group rating
			  if total_posts > 10
			}

		`)

		if transpileError != nil {
			t.Fatalf("Driver %s: Transpile error: %s", test.Driver, transpileError.Message)
		}

		AssertStringsEqual(t, test.Driver, actions[0], test.ExpectedGetUsersWithEmailAndProfile)
		AssertStringsEqual(t, test.Driver, actions[1], test.ExpectedGetUsersWithProfileAvatar)
		AssertStringsEqual(t, test.Driver, actions[2], test.ExpectedGetUsersComplex)
		AssertStringsEqual(t, test.Driver, actions[3], test.ExpectedGetUsersGroupAge)
		AssertStringsEqual(t, test.Driver, actions[4], test.ExpectedGetPostsGroupRating)
	}
}
