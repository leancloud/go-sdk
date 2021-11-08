package leancloud

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const Version = "0.3.1"

type Client struct {
	serverURL     string
	appID         string
	appKey        string
	masterKey     string
	production    string
	requestLogger *log.Logger
	Users         Users
	Files         Files
	Roles         Roles
}

type ClientOptions struct {
	AppID      string
	AppKey     string
	MasterKey  string
	ServerURL  string
	Production string
}

// NewClient constructs a client from parameters
func NewClient(options *ClientOptions) *Client {
	client := &Client{
		appID:      options.AppID,
		appKey:     options.AppKey,
		masterKey:  options.MasterKey,
		serverURL:  options.ServerURL,
		production: options.Production,
	}

	if !strings.HasSuffix(options.AppID, "MdYXbMMI") {
		if client.serverURL == "" {
			panic(fmt.Errorf("please set API's serverURL"))
		}
	}

	_, debugEnabled := os.LookupEnv("LEANCLOUD_DEBUG")

	if debugEnabled {
		client.requestLogger = log.New(os.Stdout, "", log.LstdFlags)
	}

	client.Users.c = client
	client.Files.c = client
	client.Roles.c = client
	return client
}

// NewEnvClient constructs a client from environment variables
func NewEnvClient() *Client {
	options := &ClientOptions{
		AppID:     os.Getenv("LEANCLOUD_APP_ID"),
		AppKey:    os.Getenv("LEANCLOUD_APP_KEY"),
		MasterKey: os.Getenv("LEANCLOUD_APP_MASTER_KEY"),
		ServerURL: os.Getenv("LEANCLOUD_API_SERVER"),
	}

	if appEnv := os.Getenv("LEANCLOUD_APP_ENV"); appEnv == "production" {
		options.Production = "1"
	} else if appEnv == "stage" {
		options.Production = "0"
	} else { // probably on local machine
		if os.Getenv("LEAN_CLI_HAVE_STAGING") == "true" {
			options.Production = "0"
		} else { // free trial instance only
			options.Production = "1"
		}
	}

	return NewClient(options)
}

// SetProduction sets the production environment
func (client *Client) SetProduction(production bool) {
	if production {
		client.production = "1"
	} else {
		client.production = "0"
	}
}

// Class constructs a reference of Class
func (client *Client) Class(name string) *Class {
	return &Class{
		c:    client,
		Name: name,
	}
}

// File construct a new reference to a _File object by given objectId
func (client *Client) File(id string) *FileRef {
	return &FileRef{
		c:     client,
		class: "files",
		ID:    id,
	}
}
