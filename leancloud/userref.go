package leancloud

import (
	"errors"
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

func (ref *UserRef) Get(authOptions ...AuthOption) (*User, error) {
	user := new(User)
	if err := objectGet(ref, user, authOptions...); err != nil {
		return nil, err
	}

	return user, nil
}

func (ref *UserRef) Set(field string, value interface{}, authOptions ...AuthOption) error {
	if ref.ID == "" {
		return errors.New("no reference to user")
	}

	if err := objectSet(ref, field, value, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *UserRef) Update(data map[string]interface{}, authOptions ...AuthOption) error {
	if ref.ID == "" {
		return errors.New("no reference to user")
	}

	if err := objectUpdate(ref, data, authOptions...); err != nil {
		return err
	}

	return nil
}

func (ref *UserRef) UpdateWithQuery(data map[string]interface{}, query *UserQuery, authOptions ...AuthOption) error {
	// TODO
	return nil
}

func (ref *UserRef) Destroy(authOptions ...AuthOption) error {
	if ref.ID == "" {
		return errors.New("no reference to user")
	}

	if err := objectDestroy(ref, authOptions...); err != nil {
		return err
	}

	return nil
}
