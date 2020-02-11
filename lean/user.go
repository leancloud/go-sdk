package lean

type UserInterface interface {
	Object
	getSessionToken() string
}

type User struct {
	ObjectMeta

	Email               String  `json:"email,omitempty"`
	EmailVerified       Boolean `json:"emailVerified,omitempty"`
	MobilePhoneNumber   String  `json:"mobilePhoneNumber,omitempty"`
	MobilePhoneVerified Boolean `json:"mobilePhoneVerified,omitempty"`
	SessionToken        String  `json:"sessionToken,omitempty"`
	Username            String  `json:"username,omitempty"`
	Password            String  `json:"password,omitempty"`
}

func (client *Client) Become(sessionToken string, user UserInterface) error {
	resp, err := client.request(ServiceAPI, methodGet, "/1.1/users/me", nil, UseSessionToken(sessionToken))

	if err != nil {
		return err
	}

	result := make(objectResponse)

	err = resp.JSON(&result)

	if err != nil {
		return err
	}

	return decodeObject(result, user)
}

func (*User) ClassName() string {
	return "_User"
}

func (user *User) getSessionToken() string {
	return user.SessionToken.Get()
}

func (client *Client) Login() {

}
