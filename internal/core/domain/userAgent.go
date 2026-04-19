package domain

type UserAgent struct {
	ID          int
	FullName    string
	Version     int
	PhoneNumber *string
}

func NewUserAgent(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) UserAgent {
	return UserAgent{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserAgentInitialized(
	fullName string,
	phoneNumber *string,
) UserAgent {
	return UserAgent{
		ID:          UninitializedID,
		Version:     UninitializedVersion,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}
