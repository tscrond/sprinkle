package http

import (
	"crypto/tls"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/tscrond/sprinkle/internal/db"
)

type PVERequest struct {
	Credentials *db.Credentials
	*http.Client
	*http.Request
}

func NewPVERequest(creds *db.Credentials) *PVERequest {
	return &PVERequest{Credentials: creds}
}

func (r *PVERequest) New(requestMethod, requestPath string, requestBody io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(requestMethod, requestPath, requestBody)
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
		log.Println("Initializing client transport...")
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
