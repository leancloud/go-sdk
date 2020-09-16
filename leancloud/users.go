package leancloud

import (
	"fmt"
)

type Users struct {
	c *Client
}

func (ref *Users) LogIn(username, password string) (*User, error) {
	path := fmt.Sprint("/1.1/login")
	options := ref.c.GetRequestOptions()
	options.JSON = map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := ref.c.Request(ServiceAPI, MethodPost, path, options)
	if err != nil {
		return nil, err
	}

	objectID, sessionToken, createdAt, updatedAt, respJSON, err := extracMetadata(resp.Bytes())
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

func (ref *Users) SignUp(username, password string) (*User, error) {
	user := new(User)
	reqJSON := map[string]string{
		"username": username,
		"password": password,
	}
	if err := objectCreate(ref, reqJSON, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) NewUserQuery() *UserQuery {
	return &UserQuery{}
}

func (ref *Users) Become(sessionToken string) (*User, error) {
	resp, err := ref.c.Request(ServiceAPI, MethodGet, "/1.1/users/me", ref.c.GetRequestOptions(), UseSessionToken(sessionToken))
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
