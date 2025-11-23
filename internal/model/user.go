package model

type User struct {
	ID       string 
	UserName     string 
	IsActive bool  
	TeamID   string  
}

func (u *User) AvailibleToReview() bool {
	return u.IsActive
}
