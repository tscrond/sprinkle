package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	pveLogin.Flags().Bool("envformat", false, "If flag is set, print Linux-exportable env var format")

	// Add pveLogin command to the root command
	rootCmd.AddCommand(pveLogin)
}

var pveLogin = &cobra.Command{
	Use:   "pvelogin",
	Short: "Login to PVE",
	Run: func(cmd *cobra.Command, args []string) {
		apiURL, _ := cmd.Flags().GetString("api-url")
		username, _ := cmd.Flags().GetString("username")
		envFormat, _ := cmd.Flags().GetBool("envformat")
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

		resp, err := http.Post(apiURL+"/api2/json/access/ticket", "application/json", bytes.NewBuffer(jsonData))
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
		data := result["data"].(map[string]interface{})
		csrfToken := data["CSRFPreventionToken"].(string)
		ticket := data["ticket"].(string)

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

var testConn = &cobra.Command{
	Use:   "testconn",
	Short: "Create a Proxmox resource",
	Run: func(cmd *cobra.Command, args []string) {
		// Access PersistentFlags
		apiURL, _ := cmd.Flags().GetString("api-url")
		// csrfToken, _ := cmd.Flags().GetString("csrf-token")
		// username, _ := cmd.Flags().GetString("username")
		// ticket, _ := cmd.Flags().GetString("ticket")

		csrfToken := os.Getenv("PVE_CSRFTOKEN")
		ticket := os.Getenv("PVE_TICKET")

		if csrfToken == "" {
			fmt.Println("WARNING: POST/DELETE requests will not be possible to fulfill because of no CSRF Token")
		}

		if ticket == "" {
			log.Fatalln("cannot complete request - no pve ticket")
		}

		// Create an HTTP request
		req, err := http.NewRequest("GET", apiURL+"/api2/json", nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		// Add the PVEAuthCookie to the request
		req.AddCookie(&http.Cookie{
			Name:  "PVEAuthCookie",
			Value: ticket,
		})

		// Optionally, add the CSRF token for write operations
		req.Header.Set("CSRFPreventionToken", csrfToken)

		// Send the request using an HTTP client
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// Read and print the response
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v\n", err)
			return
		}

		fmt.Println("HTTP Connectivity Status: ", resp.Status)
		fmt.Println("Response:", string(body))
	},
}
