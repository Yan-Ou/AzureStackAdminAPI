package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type ValueResponse struct {
	Value []AzureStackApp `json:"value"`
}
type AzureStackApp struct {
	ObjectID               string              `json:"objectId"`
	AppID                  string              `json:"appId"`
	AppRoleAssignments     []appRoleAssignment `json:"appRoleAssignments"`
	OAuth2PermissionGrants []interface{}       `json:"oAuth2PermissionGrants"`
	Tags                   []interface{}       `json:"tags"`
}

type appRoleAssignment struct {
	Resource string `json:"resource"`
	Client   string `json:"client"`
	RoleID   string `json:"roleId"`
}

func getRegisteredApps(token AdminApiToken, client *http.Client) ([]AzureStackApp, error) {
	var apps []AzureStackApp
	req, err := http.NewRequest("GET", config.userARMEBaseURL+"/applicationRegistrations?api-version=2014-04-01-preview", nil)
	if err != nil {
		return apps, err
	}
	req.Header.Set("Authorization", token.TokenType+" "+token.AccessToken)
	req.Header.Set("Accept", "application/json")

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return apps, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		var resStruct ValueResponse
		json.Unmarshal(bodyBytes, &resStruct)
		return resStruct.Value, nil
	}

	return apps, errors.New("Server returns: " + strconv.Itoa(res.StatusCode))
}

func registerAppHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	token, err := getToken(client, false)

	if err != nil {
		log.Fatal(err)
	}

	var apps []AzureStackApp
	apps, err = getRegisteredApps(token, client)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(apps)
}
