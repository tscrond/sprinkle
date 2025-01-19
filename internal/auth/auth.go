package auth

import (
	"errors"
	"fmt"

	"github.com/tscrond/sprinkle/internal/db"
	"gorm.io/gorm"
)

type AuthService struct {
	Db         *db.ResourceRepository
	AuthConfig []AuthConfig
}

func NewAuthService(db *db.ResourceRepository, authConfig []AuthConfig) *AuthService {
	return &AuthService{
		Db:         db,
		AuthConfig: authConfig,
	}
}


// 1. try getting creds from DB
// 2. if not - ask user for credentials
func (a *AuthService) GetCredentials(targetNode, apiUrl string) (*db.Credentials, error) {
	var creds *db.Credentials
	creds, err := a.Db.GetCredentials(targetNode, apiUrl)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		creds = a.ReadCredentialsFromUser(targetNode, apiUrl)
	}

	return creds, nil
}

func (a *AuthService) ReadCredentialsFromUser(targetNode, apiUrl string) *db.Credentials {
	var username, password string

	fmt.Print("Enter username:")
	fmt.Scanln(&username)
	fmt.Print("Enter password:")
	fmt.Scanln(&password)

	creds := &db.Credentials{
		Username:   username,
		Password:   password,
		TargetNode: targetNode,
		ApiUrl:     apiUrl,
	}

	return creds
}
