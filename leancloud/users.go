package leancloud

import (
	"encoding/json"
	"fmt"
	"time"
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

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, err
	}

	objectID, ok := respJSON["objectId"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to parse objectId from response")
	}

	createdAt, err := time.Parse(time.RFC3339, respJSON["createdAt"].(string))
	if err != nil {
		return nil, fmt.Errorf("unable to parse createdAt from response %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, respJSON["updatedAt"].(string))
	if err != nil {
		return nil, fmt.Errorf("unable to parse updatedAt from response %w", err)
	}

	sessionToken, ok := respJSON["sessionToken"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to parse sessionToken from response")
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
	path := fmt.Sprint("/1.1/users")
	options := r.c.getRequestOptions()
	options.JSON = map[string]string{
		"username": username,
		"password": password,
	}

	resp, err := r.c.request(ServiceAPI, methodPost, path, options)
	if err != nil {
		return nil, err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, err
	}

	objectID, ok := respJSON["objectId"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to parse objectId from response")
	}

	createdAt, err := time.Parse(time.RFC3339, respJSON["createdAt"].(string))
	if err != nil {
		return nil, fmt.Errorf("unable to parse createdAt from response %w", err)
	}

	sessionToken, ok := respJSON["sessionToken"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to parse sessionToken from response")
	}

	return &User{
		sessionToken: sessionToken,
		Object: Object{
			ID:        objectID,
			CreatedAt: createdAt,
		},
	}, nil
}

func (c *Client) NewUserQuery() *UserQuery {
	return &UserQuery{}
}

func (ref *Users) Become(sessionToken string) (*User, error) {
	resp, err := ref.c.request(ServiceAPI, methodPost, "/1.1/users/me", nil, UseSessionToken(sessionToken))
	if err != nil {
		return nil, err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, err
	}

	objectID, ok := respJSON["objectId"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to parse objectId from response")
	}

	createdAt, err := time.Parse(time.RFC3339, respJSON["createdAt"].(string))
	if err != nil {
		return nil, fmt.Errorf("unable to parse createdAt from response %w", err)
	}

	updatedAt, err := time.Parse(time.RFC3339, respJSON["updatedAt"].(string))
	if err != nil {
		return nil, fmt.Errorf("unable to parse updatedAt from response %w", err)
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
