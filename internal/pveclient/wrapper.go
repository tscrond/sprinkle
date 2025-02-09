package pveclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tscrond/sprinkle/internal/db"
)

type PVEClient struct {
	Credentials *db.Credentials
	*http.Client
	*http.Request
}

func NewPVEClient(creds *db.Credentials, client *http.Client) *PVEClient {
	return &PVEClient{Credentials: creds, Client: client}
}

func (r *PVEClient) NewRequest(requestMethod, requestPath string, requestBody io.Reader) (*http.Response, error) {
	fullApiUrl := fmt.Sprintf("https://%s%s", r.Credentials.ApiUrl, requestPath)

	req, err := http.NewRequest(requestMethod, fullApiUrl, requestBody)
	if err != nil {
		log.Printf("Failed to create HTTP request: %v", err)
		return nil, err
	}

	if r.Credentials.PVETicket == "" || r.Credentials.TargetNode == "" {
		log.Println("Missing authentication data: CSRF token or PVE ticket")
		return nil, errors.New("missing_auth_data")
	}

	transport, ok := r.Client.Transport.(*http.Transport)
	if !ok || transport == nil {
		// log.Println("Initializing client transport...")
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	r.Client.Transport = transport

	req.AddCookie(&http.Cookie{
		Name:  "PVEAuthCookie",
		Value: r.Credentials.PVETicket,
	})
	req.Header.Set("CSRFPreventionToken", r.Credentials.CsrfToken)

	resp, err := r.Client.Do(req)
	if err != nil {
		log.Printf("Request failed: %v", err)
		return nil, err
	}

	return resp, err
}
