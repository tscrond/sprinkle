package auth

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/tscrond/sprinkle/lib"
)

func init() {

}

var testConn = &cobra.Command{
	Use:   "testconn",
	Short: "Check connectivity with proxmox node",
	Run: func(cmd *cobra.Command, args []string) {

		apiURL, _ := cmd.Flags().GetString("api-url")

		// Create an HTTP request
		req, err := http.NewRequest("GET", "https://"+apiURL+"/api2/json", nil)
		if err != nil {
			fmt.Printf("Error creating request: %v\n", err)
			return
		}

		client, req, err := lib.ConfigureAuth(&http.Client{}, req)
		if err != nil {
			log.Fatalln("Error authenticating with PVE API")
		}

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
