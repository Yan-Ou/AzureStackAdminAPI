package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type NewDirectoryPayload struct {
	Location   string `json:"location"`
	Properties struct {
		TenantID string `json:"tenantId"`
	} `json:"properties"`
}

func addNewGuestDirectory(token AdminApiToken, guestTenantDirectory string, guestTenantID string, client *http.Client) error {
	url := config.adminARMEBaseURL + "/subscriptions/" +
		config.azureStackDefaultSubID + "/resourcegroups/" +
		config.azureStackResourceGroup + "/providers/Microsoft.Subscriptions.Admin/directoryTenants/" +
		guestTenantDirectory + "?api-version=2015-11-01"
	var payload NewDirectoryPayload
	payload.Location = config.azureStackRegion
	payload.Properties.TenantID = guestTenantID
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	var req *http.Request
	req, err = http.NewRequest("PUT", url, bytes.NewBuffer(b))
	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusCreated {
		return nil
	}
	return err
}

func onBoardHandler(w http.ResponseWriter, r *http.Request) {
	guestTenantID := r.FormValue("guestTenantId")               // f71f054e-ebf4-4526-8678-4d85999db2ab
	guestTenantDirectory := r.FormValue("guestTenantDirectory") // azurestacktestmax.onmicrosoft.com

	if guestTenantID == "" || guestTenantDirectory == "" {
		http.Error(w, "Missing paramerts", 500)
		return
	}

	client := &http.Client{}
	token, err := getToken(client, true)
	if err != nil {
		log.Fatal(err)
	}

	if err = addNewGuestDirectory(token, guestTenantID, guestTenantDirectory, client); err != nil {
		log.Fatal(err)
	} else {
		fmt.Print("success")
	}
}
