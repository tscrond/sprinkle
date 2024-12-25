package lib

import (
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"os"
)

func ConfigureAuth(client *http.Client, req *http.Request) (*http.Client, *http.Request, error) {

	csrfToken := os.Getenv("PVE_CSRFTOKEN")
	ticket := os.Getenv("PVE_TICKET")

	if csrfToken == "" {
		log.Println("WARNING: POST/DELETE requests will not be possible to fulfill because of no CSRF Token")
	}

	if ticket == "" {
		log.Println("cannot complete request - no pve ticket")
		return nil, nil, errors.New("no_pve_ticket")
	}

	transport, ok := client.Transport.(*http.Transport)
	if !ok || transport == nil {
		// Create a new http.Transport if it doesn't exist
		transport = &http.Transport{}
	}

	transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client.Transport = transport

	// Add the PVEAuthCookie to the request
	req.AddCookie(&http.Cookie{
		Name:  "PVEAuthCookie",
		Value: ticket,
	})

	// Optionally, add the CSRF token for write operations
	req.Header.Set("CSRFPreventionToken", csrfToken)

	return client, req, nil
}
