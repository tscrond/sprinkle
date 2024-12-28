package lib

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"os"
)

func GetReq(apiUrlFull string) (*http.Response, error) {
	req, _ := http.NewRequest("GET", apiUrlFull, nil)
	client := &http.Client{}

	client, req, err := ConfigureAuth(client, req)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
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

func assignIDToNode(apiUrl string, currentID int) int {
	err := checkIDExists(apiUrl, currentID)
	if err == nil {
		id := currentID
		fmt.Println("Assigning VMID: ", id)
		return id
	}

	if !errors.Is(err, ErrIDExists) {
		log.Println(err)
		return -1
	}

	newID := rand.IntN(9999999)

	return assignIDToNode(apiUrl, newID)
}

func checkIDExists(apiUrl string, id int) error {
	nodes := discoverNodes(apiUrl)

	type machine struct {
		VMID int `json:"vmid"`
	}

	for _, node := range nodes {
		apiUrlFull := "https://" + apiUrl + "/api2/json/nodes/" + node

		var lxcs struct {
			Data []machine `json:"data"`
		}
		var vms struct {
			Data []machine `json:"data"`
		}

		respLxc, err := GetReq(apiUrlFull + "/lxc")
		if err != nil {
			return err
		}
		if err := json.NewDecoder(respLxc.Body).Decode(&lxcs); err != nil {
			return err
		}
		defer respLxc.Body.Close()

		respVm, err := GetReq(apiUrlFull + "/qemu")
		if err != nil {
			return err
		}
		if err := json.NewDecoder(respVm.Body).Decode(&vms); err != nil {
			return err
		}
		defer respVm.Body.Close()

		allVms := append(lxcs.Data, vms.Data...)
		for _, vm := range allVms {
			if vm.VMID == id {
				return errors.New("id_exists")
			}
		}
	}

	return nil
}

func discoverNodes(apiUrl string) []string {
	req, _ := http.NewRequest("GET", "https://"+apiUrl+"/api2/json/nodes", nil)
	client := &http.Client{}

	client, req, err := ConfigureAuth(client, req)
	if err != nil {
		return []string{}
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	type pvenode struct {
		Node string `json:"node"`
	}

	var result struct {
		Data []pvenode `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return []string{}
	}

	var nodes []string
	for _, node := range result.Data {
		nodes = append(nodes, node.Node)
	}

	return nodes
}
