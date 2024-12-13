package db

import (
	"database/sql"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type userDBImpl struct {
	DB *sql.DB
}

func NewUserDBImpl(db *sql.DB) *userDBImpl {
	return &userDBImpl{DB: db}
}

func (u *userDBImpl) GetUserByEmail(email string) (*user, error) {
	query := "SELECT id, email, password, first_name, last_name, created_at FROM users WHERE email = $1"
	var user user
	err := u.DB.QueryRow(query, email).Scan(&user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName, &user.CreatedAt)
	if err != nil {
		log.Println("error getting user by email:", err)
		return nil, err
	}
	return &user, nil
}

func (u *userDBImpl) PasswordCheck(hashedPassword string, plainPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	if err != nil {
		log.Println("error validating user:", err)
		return err
	}
	return nil
}
