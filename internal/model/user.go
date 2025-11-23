package model

type User struct {
	ID       string
	Username string
	IsActive bool
	TeamName string
}

func (u *User) AvailibleToReview() bool {
	return u.IsActive
}
