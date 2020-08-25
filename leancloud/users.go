package leancloud

import (
	"fmt"
)

type Users struct {
	c *Client
}

func (r *Users) LogIn(username, password string) (*User, error) {
	path := fmt.Sprint("/1.1/login")
	options := r.c.getRequestOptions()
	options.JSON = map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := r.c.request(ServiceAPI, methodPost, path, options)
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

func (r *Users) SignUp(username, password string) (*User, error) {
	user := new(User)
	reqJSON := map[string]string{
		"username": username,
		"password": password,
	}
	if err := objectCreate(r, reqJSON, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) NewUserQuery() *UserQuery {
	return &UserQuery{}
}

func (ref *Users) Become(sessionToken string) (*User, error) {
	resp, err := ref.c.request(ServiceAPI, methodPost, "/1.1/users/me", nil, UseSessionToken(sessionToken))
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
