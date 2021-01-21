package leancloud

import (
	"log"
	"os"
)

const Version = "0.1.0"

type Client struct {
	serverURL     string
	appID         string
	appKey        string
	masterKey     string
	requestLogger *log.Logger
	Users         Users
	Files         Files
	Roles         Roles
}

func (client *Client) SetServerURL(url string) {
	client.serverURL = url
}

// NewClient constructs a client from parameters
func NewClient(appID, appKey, masterKey string) *Client {
	client := &Client{
		appID:     appID,
		appKey:    appKey,
		masterKey: masterKey,
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
	client := NewClient(os.Getenv("LEANCLOUD_APP_ID"), os.Getenv("LEANCLOUD_APP_KEY"), os.Getenv("LEANCLOUD_APP_MASTER_KEY"))
	client.SetServerURL(os.Getenv("LEANCLOUD_API_SERVER"))
	return client
}

// Class constrcuts a reference of Class
func (client *Client) Class(name string) *Class {
	return &Class{
		c:    client,
		Name: name,
	}
}

// File construct an new reference to a _File object by given objectId
func (client *Client) File(id string) *FileRef {
	return &FileRef{
		c:     client,
		class: "files",
		ID:    id,
	}
}
