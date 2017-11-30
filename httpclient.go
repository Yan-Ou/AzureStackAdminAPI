package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const directoryID = "5c37fe2b-b77a-4388-a815-201eb361be16"
const clientSecret = "x5lXipWyEJhXRYV4sCk6Y1+zKzhifa8du7DInNS1lDE="
const username = "ublradmin@umbrellarstack2.onmicrosoft.com"
const password = "AzureStack!!"

//const ArmEndpoit = "https://adminmanagement.asdk2.umbrellar.io"
const resource = "https://adminmanagement.umbrellarstack2.onmicrosoft.com/c3da154f-8979-424d-a256-3fbd4e9488e2"
const aadURL = "https://login.microsoftonline.com/" + directoryID + "/oauth2/token?api-version=1.0"
const clientID = "dda51c56-a86e-452f-a91f-cef7b0785526"

//const defaultSubID = "f0dcdd97-386f-4af2-bbbc-26152c89ad07"

type authToken struct {
	AccessToken string `json:"access_token"`
}

type Subscriptions struct {
	Subscription []azSub `json:"value"`
}

type azSub struct {
	SubID    string `json:"subscriptionId"`
	Name     string `json:"displayName"`
	Owner    string `json:"owner"`
	TenantID string `json:"tenantId"`
}

func getToken() string {
	form := url.Values{
		"grant_type":    {"password"},
		"scope":         {"openid"},
		"resource":      {resource},
		"client_id":     {clientID},
		"client_secret": {clientSecret},
		"username":      {username},
		"password":      {password}}
	resp, err := http.PostForm(aadURL, form)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var accessToken authToken
	err = json.Unmarshal(body, &accessToken)
	if err != nil {
		log.Fatal(err)
	}

	return accessToken.AccessToken
}

func main() {
	aToken := getToken()
	fmt.Print(aToken)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://adminmanagement.asdk2.umbrellar.io/subscriptions/f0dcdd97-386f-4af2-bbbc-26152c89ad07/providers/Microsoft.Subscriptions.Admin/subscriptions?api-version=2015-11-01", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+aToken)
	resp, _ := client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)
	var sub Subscriptions
	err = json.Unmarshal(body, &sub)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%+v", sub)
}
