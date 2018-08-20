package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type EndpointMetaData struct {
	GalleryEndpoint string `json:"galleryEndpoint"`
	GraphEndpoint   string `json:"graphEndpoint"`
	PortalEndpoint  string `json:"portalEndpoint"`
	Authentication  struct {
		LoginEndpoint string   `json:"loginEndpoint"`
		Audiences     []string `json:"audiences"`
	} `json:"authentication"`
}

type AdminApiToken struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    string `json:"expires_in"`
	ExtExpiresIn string `json:"ext_expires_in"`
	ExpiresOn    string `json:"expires_on"`
	NotBefore    string `json:"not_before"`
	Resource     string `json:"resource"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

func getToken(client *http.Client, admin bool) (AdminApiToken, error) {
	var token AdminApiToken
	audienceAuthenticationEndpoint, err := getAudienceAuthenticationEndpoint(client, admin)
	if err != nil {
		return token, err
	}

	token, err = getApiAdminARMToken(audienceAuthenticationEndpoint, client)
	return token, err
}

func getAudienceAuthenticationEndpoint(client *http.Client, admin bool) (string, error) {
	var url string
	if admin {
		url = config.adminARMEBaseURL
	} else {
		url = config.userARMEBaseURL
	}
	req, err := http.NewRequest("GET", url+"/metadata/endpoints?api-version=2015-01-01", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json")
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		var raw EndpointMetaData
		json.Unmarshal(bodyBytes, &raw)
		return raw.Authentication.Audiences[0], nil
	}

	return "", errors.New("Server returns: " + strconv.Itoa(res.StatusCode))

}

func getApiAdminARMToken(audienceAuthenticationEndpoint string, client *http.Client) (AdminApiToken, error) {
	var token AdminApiToken
	form := url.Values{}
	form.Add("resource", audienceAuthenticationEndpoint)
	form.Add("client_id", config.azureStackClientAppID)
	form.Add("client_secret", config.azureStackClientAppSecret)
	form.Add("grant_type", "password")
	form.Add("username", config.oAuthUsername)
	form.Add("password", config.oAuthPassword)
	form.Add("scope", "openid")
	req, err := http.NewRequest("POST", "https://login.microsoftonline.com/"+config.oAuthDomain+"/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return token, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return token, err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(bodyBytes, &token)
		return token, nil
	}

	return token, errors.New("Server returns: " + strconv.Itoa(res.StatusCode))
}
