package leancloud

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func (ref *Users) LogInByMobilePhoneNumber(number, smsCode string) (*User, error) {
	path := "/1.1/login"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"mobilePhoneNumber": number,
		"smsCode":           smsCode,
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

func (ref *Users) SignUpByMobilePhone(number, smsCode string) (*User, error) {
	path := "/1.1/usersByMobilePhone"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"mobilePhoneNumber": number,
		"smsCode":           smsCode,
	}

	resp, err := ref.c.request(ServiceAPI, MethodPost, path, options)
	if err != nil {
		return nil, err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, err
	}

	decodedUser, err := decodeUser(respJSON)
	if err != nil {
		return nil, err
	}

	return decodedUser, nil
}

func (ref *Users) ResetPasswordBySMSCode(number, smsCode, password string, authOptions ...AuthOption) error {
	path := "/1.1/resetPasswordBySmsCode/"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"password":          password,
		"mobilePhoneNumber": number,
	}

	_, err := ref.c.request(ServiceAPI, MethodPost, fmt.Sprint(path, smsCode), options, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) NewUserQuery() *UserQuery {
	return &UserQuery{}
}

func (ref *Users) Become(sessionToken string) (*User, error) {
	resp, err := ref.c.request(ServiceAPI, MethodGet, "/1.1/users/me", ref.c.getRequestOptions(), UseSessionToken(sessionToken))
	if err != nil {
		return nil, err
	}

	respJSON := make(map[string]interface{})
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		return nil, err
	}

	return decodeUser(respJSON)
}

func (ref *Users) RequestEmailVerify(email string, authOptions ...AuthOption) error {
	path := "/1.1/requestEmailVerify"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"email": email,
	}

	_, err := ref.c.request(ServiceAPI, MethodPost, path, options, authOptions...)
	if err != nil {
		return err
	}

	return nil
}

func (ref *Users) RequestMobilePhoneVerify(number string, authOptions ...AuthOption) error {
	path := "/1.1/requestMobilePhoneVerify"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"mobilePhoneNumber": number,
	}

	resp, err := ref.c.request(ServiceAPI, MethodPost, path, options, authOptions...)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(resp.Bytes()))
	}

	return nil
}

func (ref *Users) RequestPasswordReset(email string, authOptions ...AuthOption) error {
	path := "/1.1/requestPasswordReset"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"email": email,
	}

	resp, err := ref.c.request(ServiceAPI, MethodPost, path, options, authOptions...)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(resp.Bytes()))
	}

	return nil
}

func (ref *Users) RequestPasswordResetBySMSCode(number string, authOptions ...AuthOption) error {
	path := "/1.1/requestPasswordResetBySmsCode"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"mobilePhoneNumber": number,
	}

	resp, err := ref.c.request(ServiceAPI, MethodPost, path, options, authOptions...)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(resp.Bytes()))
	}

	return nil
}

func (ref *Users) RequestLoginSMSCode(number string, authOptions ...AuthOption) error {
	path := "/1.1/requestLoginSmsCode"
	options := ref.c.getRequestOptions()
	options.JSON = map[string]string{
		"mobilePhoneNumber": number,
	}

	resp, err := ref.c.request(ServiceAPI, MethodPost, path, options, authOptions...)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s", string(resp.Bytes()))
	}

	return nil
}
