package leancloud

import (
	"log"
	"os"
)

const Version = "0.1.0"

type Client struct {
	region        Region
	appID         string
	appKey        string
	masterKey     string
	requestLogger *log.Logger
	Users         Users
	Files         Files
}

func NewClient(region, appID, appKey, masterKey string) *Client {
	client := &Client{
		region:    NewRegionFromString(region),
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
	return client
}
func NewEnvClient() *Client {
	return NewClient(os.Getenv("LEANCLOUD_REGION"),
		os.Getenv("LEANCLOUD_APP_ID"),
		os.Getenv("LEANCLOUD_APP_KEY"),
		os.Getenv("LEANCLOUD_APP_MASTER_KEY"))
}
