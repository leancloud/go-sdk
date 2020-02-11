package lean

import (
	"fmt"
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
}

func NewEnvClient() *Client {
	client := &Client{
		region:    NewRegionFromString(os.Getenv("LEANCLOUD_REGION")),
		appID:     os.Getenv("LEANCLOUD_APP_ID"),
		appKey:    os.Getenv("LEANCLOUD_APP_KEY"),
		masterKey: os.Getenv("LEANCLOUD_APP_MASTER_KEY"),
	}

	_, foundDebugEnv := os.LookupEnv("LEANCLOUD_DEBUG")

	if foundDebugEnv {
		client.requestLogger = log.New(os.Stdout, "", log.LstdFlags)
	}

	return client
}

func (client *Client) Save(object Object, authOptions ...AuthOption) error {
	requestBody := map[string]interface{}{}

	err := encodeObject(object, requestBody)

	if err != nil {
		return err
	}

	method := methodPost
	path := fmt.Sprint("/1.1/classes/", object.ClassName())

	if object.getObjectMeta().ObjectID != "" {
		method = methodPut
		path = fmt.Sprint(path, "/", object.getObjectMeta().ObjectID)
	}

	options := client.getRequestOptions()

	options.JSON = requestBody

	resp, err := client.request(ServiceAPI, method, path, options, authOptions...)

	if err != nil {
		return err
	}

	// result := &createObjectResponse{}

	// err = resp.JSON(result)
	//
	// if err != nil {
	// 	return err
	// }

	err = mergeToObject(resp.Bytes(), object)

	if err != nil {
		return err
	}

	return nil
}

func (client *Client) Destroy(object Object, authOptions ...AuthOption) error {
	classPath := object.ClassName()

	if classPath == "_User" {
		classPath = "users"
	} else {
		classPath = "classes/users"
	}

	path := fmt.Sprint("/1.1/", classPath, "/", object.getObjectMeta().ObjectID)

	_, err := client.request(ServiceAPI, methodDelete, path, nil, authOptions...)

	return err
}

func mergeDataFromServer(object Object, resp *createObjectResponse) {
	meta := object.getObjectMeta()

	meta.ObjectID = resp.ObjectID
	meta.CreatedAt = resp.CreatedAt
}
