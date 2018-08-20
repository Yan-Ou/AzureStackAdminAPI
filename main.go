package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Config struct {
	adminARMEBaseURL          string `yaml:"adminARMEBaseURL"`
	userARMEBaseURL           string `yaml:"userARMEBaseURL"`
	azureStackClientAppID     string `yaml:"azureStackClientAppID"`
	azureStackClientAppSecret string `yaml:"azureStackClientAppSecret"`
	oAuthDomain               string `yaml:"oAuthDomain"`
	oAuthUsername             string `yaml:"oAuthUsername"`
	oAuthPassword             string `yaml:"oAuthPassword"`
	azureStackDefaultSubID    string `yaml:"azureStackDefaultSubID"`
	azureStackResourceGroup   string `yaml:"azureStackResourceGroup"`
	azureStackRegion          string `yaml:"azureStackRegion"`
}

func (c *Config) Load(env string) error {
	if env == "development" {
		c.adminARMEBaseURL = "https://adminmanagement.asdk1.umbrellar.io"
		c.userARMEBaseURL = "https://management.asdk1.umbrellar.io"
		c.azureStackClientAppID = "c664b98b-25c2-4430-9469-1abdb4379510"
		c.azureStackClientAppSecret = "NPwoUXcEyw4SEdSt0otawowiY/6OockJ+joEa+VXBzI="
		c.oAuthDomain = "ublrsandbox.onmicrosoft.com"
		c.oAuthUsername = "portal.agent@ublrsandbox.onmicrosoft.com"
		c.oAuthPassword = "XkqVpFkuL8l@vFwN90P&0t8z"
		c.azureStackDefaultSubID = "180974cf-aaa5-47e2-9bc5-9949f9c9ac70"
		c.azureStackResourceGroup = "TenantDirectory"
		c.azureStackRegion = "asdk1"

	} else {
		c.adminARMEBaseURL = os.Getenv("adminARMEBaseURL")
		c.userARMEBaseURL = os.Getenv("userARMEBaseUrl")
		c.azureStackClientAppID = os.Getenv("azureStackClientAppID")
		c.azureStackClientAppSecret = os.Getenv("azureStackClientAppSecret")
		c.oAuthDomain = os.Getenv("oAuthDomain")
		c.oAuthUsername = os.Getenv("oAuthUsername")
		c.oAuthPassword = os.Getenv("oAuthPassword")
		c.azureStackDefaultSubID = os.Getenv("azureStackDefaultSubID")
		c.azureStackResourceGroup = os.Getenv("azureStackResourceGroup")
		c.azureStackRegion = os.Getenv("azureStackRegion")
	}
	return nil
}

var config Config

func main() {
	env := os.Getenv("APP_ENV")
	if err := config.Load(env); err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/onboard-new-tenant", onBoardHandler).Methods("POST")
	r.HandleFunc("/register-apps", registerAppHandler).Methods("GET")
	log.Println("Up and Serving")
	if err := http.ListenAndServe(":9000", r); err != nil {
		log.Fatal(err)
	}
}
