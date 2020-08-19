package leancloud

import "time"

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
	// TODO
	return nil
}

func (ref *UserRef) Get(authOption ...AuthOption) (*User, error) {
	// TODO
	return nil, nil
}

func (ref *UserRef) Set(field string, value interface{}, authOption ...AuthOption) error {
	// TODO
	return nil
}

func (ref *UserRef) Update(data map[string]interface{}, authOption ...AuthOption) error {
	// TODO
	return nil
}

func (ref *UserRef) UpdateWithQuery(data map[string]interface{}, query *UserQuery, authOption ...AuthOption) error {
	// TODO
	return nil
}

func (ref *UserRef) Delete() error {
	// TODO
	return nil
}
