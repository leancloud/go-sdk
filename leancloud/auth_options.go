package leancloud

import (
	"fmt"

	"github.com/levigross/grequests"
)

type AuthOption interface {
	apply(*Client, *grequests.RequestOptions)
}

type authOption struct {
	useMasterKey bool
	sessionToken string
}

func (option *authOption) apply(client *Client, request *grequests.RequestOptions) {
	if option.useMasterKey {
		request.Headers["X-LC-Key"] = fmt.Sprint(client.masterKey, ",master")
	}

	if option.sessionToken != "" {
		request.Headers["X-LC-Session"] = option.sessionToken
	}
}

func UseMasterKey(useMasterKey bool) AuthOption {
	return &authOption{
		useMasterKey: useMasterKey,
	}
}

func UseSessionToken(sessionToken string) AuthOption {
	return &authOption{
		sessionToken: sessionToken,
	}
}

func UseUser(user *User) AuthOption {
	return &authOption{
		sessionToken: user.SessionToken(),
	}
}
