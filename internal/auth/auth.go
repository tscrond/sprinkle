package auth

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tscrond/sprinkle/internal/db"
	"gorm.io/gorm"
)

type AuthService struct {
	Db *db.ResourceRepository
}

func NewAuthService(db *db.ResourceRepository) *AuthService {
	return &AuthService{
		Db: db,
	}
}

// authenticated request to proxmox api to target node
// func (a *AuthService) ProxmoxAPIRequest(method, requestPath, targetNode string, requestBody io.Reader) (*http.Response, error) {

//		return nil, nil
//	}

// return auth creds for target node
func (a *AuthService) Authenticate(targetNode, apiUrl string) (*db.Credentials, error) {
	creds, err := a.GetCredentials(targetNode, apiUrl)
	if err != nil {
		return nil, err
	}

	creds, err = a.RetrieveTokenAndTicket(creds)
	if err != nil {
		// fmt.Println(err)
		return nil, err
	}

	if err := a.PersistCredsInDB(creds); err != nil {
		fmt.Println(err)
	}

	return creds, nil
}

func (a *AuthService) RetrieveTokenAndTicket(userProvidedCreds *db.Credentials) (*db.Credentials, error) {

	if (userProvidedCreds.CsrfToken != "" && userProvidedCreds.PVETicket != "") && !time.Now().After(userProvidedCreds.UpdatedAt.Add(2*time.Hour)) {
		return userProvidedCreds, nil
	}

	pveAuthData := map[string]string{
		"username": userProvidedCreds.Username,
		"password": userProvidedCreds.Password,
	}

	jsonData, err := json.Marshal(pveAuthData)
	if err != nil {
		fmt.Printf("Error marshalling auth data: %v\n", err)
		return nil, err
	}

	// Create a custom HTTP client with TLS verification disabled
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	// Prepare the request to the Proxmox API
	req, err := http.NewRequest("POST", "https://"+userProvidedCreds.ApiUrl+"/api2/json/access/ticket", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return nil, err
	}

	// set headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Printf("Error making request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	// Parse response JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return nil, err
	}

	// Extract CSRF token and ticket from the response
	data, ok := result["data"].(map[string]interface{})
	if !ok {
		authFailureErr := errors.New("auth_failure")
		fmt.Println("Error:", authFailureErr)
		return nil, authFailureErr
	}

	csrfToken, _ := data["CSRFPreventionToken"].(string)
	ticket, _ := data["ticket"].(string)

	// if csrfToken == "" && ticket == "" {
	// 	return nil, errors.New("auth_failure")
	// }

	creds := &db.Credentials{
		TargetNode: userProvidedCreds.TargetNode,
		ApiUrl:     userProvidedCreds.ApiUrl,
		Username:   userProvidedCreds.Username,
		Password:   userProvidedCreds.Password,
		CsrfToken:  csrfToken,
		PVETicket:  ticket,
	}

	// Return the CSRF token and ticket
	return creds, err
}

func (a *AuthService) PersistCredsInDB(creds *db.Credentials) error {
	if a.Db.CredsExist(creds.TargetNode, creds.ApiUrl) {
		return a.Db.UpdateCredentials(creds)
	}
	return a.Db.SaveCredentials(creds)
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

	fmt.Println("Authenticating with node: ", targetNode)
	// for now, shitty credential reading
	fmt.Print("Enter username:")
	fmt.Scanln(&username)
	fmt.Print("Enter password:")
	fmt.Scanln(&password)

	creds := &db.Credentials{
		Username:   username + "@pam",
		Password:   password,
		TargetNode: targetNode,
		ApiUrl:     apiUrl,
	}

	return creds
}
