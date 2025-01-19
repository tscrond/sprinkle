package auth

import "github.com/tscrond/sprinkle/internal/db"

type AuthService struct {
	Db          *db.ResourceRepository
	Credentials *Credentials
}

func NewAuthService(db *db.ResourceRepository, creds *Credentials) *AuthService {
	return &AuthService{
		Db:          db,
		Credentials: creds,
	}
}

func (ac *AuthService) GetCredentials(clusterName string) (*Credentials, error) {
	return nil, nil
}
