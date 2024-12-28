package auth

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	pveLogin.Flags().Bool("env-format", false, "If flag is set, print Linux-exportable env var format")
	pveLogin.Flags().String("username", "root@pam", "user name example: 'root@pam'")
}

var pveLogin = &cobra.Command{
	Use:   "pvelogin",
	Short: "Login to PVE",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Flags().GetString("api-url")
		username, _ := cmd.Flags().GetString("username")
		envFormat, _ := cmd.Flags().GetBool("env-format")
		password := os.Getenv("PVE_PASSWORD")

		if password == "" {
			log.Fatalln("No PVE_PASSWORD set, cannot authenticate")
		}

		authData := map[string]string{
			"username": username,
			"password": password,
		}

		jsonData, err := json.Marshal(authData)
		if err != nil {
			fmt.Printf("Error marshalling auth data: %v\n", err)
			return
		}

		// Create a custom HTTP client with TLS verification disabled
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}

		req, err := http.NewRequest("POST", "https://"+apiURL+"/api2/json/access/ticket", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}

		// Parse response JSON
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("Error parsing response: %v\n", err)
			return
		}

		// Extract token and ticket
		data, ok := result["data"].(map[string]interface{})
		if !ok {
			fmt.Println("Error: Invalid response format")
			return
		}

		csrfToken, _ := data["CSRFPreventionToken"].(string)
		ticket, _ := data["ticket"].(string)

		if envFormat {
			fmt.Println("#!/bin/bash")
			fmt.Printf("export PVE_CSRFTOKEN=\"%s\"\n", csrfToken)
			fmt.Printf("export PVE_TICKET=\"%s\"\n", ticket)
		} else {
			fmt.Println("CSRF Token:", csrfToken)
			fmt.Println("Ticket:", ticket)
		}
	},
}
