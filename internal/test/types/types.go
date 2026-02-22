package test_types

type ToTestCreateTable struct {
	Driver                string
	ExpectedUsersTable    string
	ExpectedProfilesTable string
	ExpectedPostsTable    string
}

type ToTestCreateGet struct {
	Driver                              string
	ExpectedGetUsersWithEmailAndProfile string
	ExpectedGetUsersWithProfileAvatar   string
	ExpectedGetUsersComplex             string
	ExpectedGetUsersGroupAge            string
	ExpectedGetPostsGroupRating         string
}

type ToTestCreateDelete struct {
	Driver             string
	ExpectedDeleteUser string
}

type ToTestCreateNew struct {
	Driver              string
	ExpectedNewUser     string
	ExpectedNewBulkUser string
}
