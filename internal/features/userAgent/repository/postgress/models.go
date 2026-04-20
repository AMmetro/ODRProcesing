package users_postgres_repository

type UserAgentModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}
