package db

type UserDB interface {
	GetUserByEmail(email string) (*user, error)
	PasswordCheck(hashedPassword, plainPassword string) error
}
