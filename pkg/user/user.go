package user

type User struct {
	ID       string
	Login    string
	password string
}

type UserRepo interface {
	Authorize(login, pass string) (*User, error)
	AddUser(login, pass string) error
}
