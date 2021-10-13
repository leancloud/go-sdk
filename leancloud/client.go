package leancloud

import (
	"fmt"
	"log"
	"os"
	"strings"
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

type ClientOptions struct {
	AppID     string
	AppKey    string
	MasterKey string
	ServerURL string
}

// NewClient constructs a client from parameters
func NewClient(options *ClientOptions) *Client {
	client := &Client{
		appID:     options.AppID,
		appKey:    options.AppKey,
		masterKey: options.MasterKey,
		serverURL: options.ServerURL,
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

	return NewClient(options)
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
