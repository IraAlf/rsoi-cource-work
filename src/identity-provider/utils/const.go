package utils

type Roles int

const (
	Admin Roles = iota + 1
	User
)

func (me Roles) String() string {
	return [...]string{"", "admin", "user"}[me]
}
