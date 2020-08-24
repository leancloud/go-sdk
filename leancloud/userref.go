package leancloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type UserRef struct {
	c     *Client
	class string
	ID    string
}

type signupResponse struct {
	SessionToken string    `json:"sessionToken"`
	CreatedAt    time.Time `json:"createdAt"`
	ObjectID     string    `json:"objectId"`
}

type signinResponse struct {
	SessionToken        string    `json:"sessionToken"`
	UpdatedAt           time.Time `json:"updatedAt"`
	Phone               string    `json:"phone"`
	ObjectID            string    `json:"objectId"`
	Username            string    `json:"username"`
	CreatedAt           time.Time `json:"createdAt"`
	EamilVerified       bool      `json:"emailVerified"`
	MobilePhoneVerified bool      `json:"mobilePhoneVerified"`
}

func (client *Client) User(id string) *UserRef {
	return &UserRef{
		c:     client,
		class: "users",
		ID:    id,
	}
}

func (ref *UserRef) Get(authOption ...AuthOption) (*User, error) {
	path := fmt.Sprint("/1.1/classes/users/", ref.ID)

	resp, err := ref.c.request(ServiceAPI, methodGet, path, ref.c.getRequestOptions(), authOption...)
	if err != nil {
		return nil, err
	}

	resBody := make(map[string]interface{})

	if err := json.Unmarshal(resp.Bytes(), &resBody); err != nil {
		return nil, err
	}

	createdAt, err := time.Parse(time.RFC3339, resBody["createdAt"].(string))
	if err != nil {
		return nil, err
	}
	updatedAt, err := time.Parse(time.RFC3339, resBody["updatedAt"].(string))
	if err != nil {
		return nil, err
	}

	sessionToken, ok := resBody["sessionToken"].(string)
	if !ok {
		return nil, errors.New("unable to parse sessionToken from response")
	}
	objectID, ok := resBody["objectId"].(string)
	if !ok {
		return nil, errors.New("unable to parse objectId from response")
	}

	return &User{
		sessionToken: sessionToken,
		Object: Object{
			ID:        objectID,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
			fields:    resBody,
		},
	}, nil
}

func (ref *UserRef) Set(field string, value interface{}, authOption ...AuthOption) error {
	if ref.ID == "" {
		return errors.New("no reference to user")
	}

	path := fmt.Sprint("/1.1/classes/users/", ref.ID)
	options := ref.c.getRequestOptions()
	options.JSON = encodeObject(map[string]interface{}{
		field: value,
	})

	resp, err := ref.c.request(ServiceAPI, methodPut, path, options, authOption...)
	if err != nil {
		return err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	return nil
}

func (ref *UserRef) Update(data map[string]interface{}, authOption ...AuthOption) error {
	if ref.ID == "" {
		return errors.New("no reference to user")
	}

	path := fmt.Sprint("/1.1/classes/users/", ref.ID)
	options := ref.c.getRequestOptions()
	options.JSON = encodeObject(data)

	resp, err := ref.c.request(ServiceAPI, methodPut, path, options, authOption...)
	if err != nil {
		return err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	return nil
}

func (ref *UserRef) UpdateWithQuery(data map[string]interface{}, query *UserQuery, authOption ...AuthOption) error {
	// TODO
	return nil
}

func (ref *UserRef) Delete(authOption ...AuthOption) error {
	if ref.ID == "" {
		return errors.New("no reference to user")
	}

	path := fmt.Sprint("/1.1/classes/users/", ref.ID)
	resp, err := ref.c.request(ServiceAPI, methodDelete, path, ref.c.getRequestOptions(), authOption...)
	if err != nil {
		return err
	}

	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return err
	}

	return nil
}
