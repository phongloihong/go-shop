package entity

type UserPublicProfile struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func NewUserPublicProfile(id, firstName, lastName string) *UserPublicProfile {
	profile := &UserPublicProfile{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
	}

	return profile
}
