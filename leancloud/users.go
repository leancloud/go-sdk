package leancloud

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type Users struct {
	c *Client
}

func (ref *Users) LogIn(username, password string) (*User, error) {
	path := fmt.Sprint("/1.1/login")
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := ref.c.request(ServiceAPI, MethodPost, path, options)
	if err != nil {
		return nil, err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, err
	}

	return decodeUser(respJSON)
}

func (ref *Users) SignUp(username, password string) (*User, error) {
	reqJSON := map[string]string{
		"username": username,
		"password": password,
	}
	decodedUser, err := objectCreate(ref, reqJSON)
	if err != nil {
		return nil, err
	}

	user, ok := decodedUser.(*User)
	if !ok {
		return nil, fmt.Errorf("unexpected error when parse User from response: want type *User but %v", reflect.TypeOf(decodedUser))
	}
	return user, nil
}

func (c *Client) NewUserQuery() *UserQuery {
	return &UserQuery{}
}

func (ref *Users) Become(sessionToken string) (*User, error) {
	resp, err := ref.c.request(ServiceAPI, MethodGet, "/1.1/users/me", ref.c.getRequestOptions(), UseSessionToken(sessionToken))
	if err != nil {
		return nil, err
	}

	objectID, _, createdAt, updatedAt, respJSON, err := extracMetadata(resp.Bytes())
	if err != nil {
		return nil, err
	}

	return &User{
		sessionToken: sessionToken,
		Object: Object{
			ID:        objectID,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			fields:    respJSON,
		},
	}, nil
}
