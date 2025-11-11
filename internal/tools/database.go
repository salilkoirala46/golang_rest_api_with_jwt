package tools

import "github.com/gofiber/fiber/v2/log"

type Token struct {
	Token string
}

type UserDetail struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DatabaseInterface interface {
	RegisterUser(user *UserDetail, password string) (*Token, error)
	Login(user *UserDetail, password string) (*Token, error)
	ValidateToken(token string) (*Token, error)
	SetupDatabase() error
	GetAllUsers() (*[]UserDetail, error)
	AddUser(user *UserDetail) (*[]UserDetail, error)
	UpdateUser(updatedUser UserDetail, id int) (*[]UserDetail, error)
	DeleteUser(id int) (*[]UserDetail, error)
}

func NewDatabase() (*DatabaseInterface, error) {
	// var db DatabaseInterface = &mockDB{}
	var db DatabaseInterface = &postgresDB{}
	var err error = db.SetupDatabase()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &db, nil

}
