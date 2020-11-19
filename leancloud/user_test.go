package leancloud

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

type CustomUser struct {
	Username            string `json:"username"`
	MobilePhoneVerified bool   `json:"mobilePhoneVerified"`
	EmailVerified       bool   `json:"emailVerified"`
}

func TestUserGetMap(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	userMap := user.GetMap()

	if userMap["username"].(string) != username {
		t.Fatal(fmt.Errorf("username unmatch"))
	}

	if userMap["sessionToken"].(string) != user.sessionToken {
		t.Fatal(fmt.Errorf("sessionToken unmatch"))
	}

	createdAt, ok := userMap["createdAt"].(time.Time)
	if !ok {
		t.Log(reflect.TypeOf(userMap["createdAt"]))
		t.Fatal(fmt.Errorf("unable to parse createdAt from fields"))
	}

	if createdAt.Unix() != user.CreatedAt.Unix() {
		t.Fatal(fmt.Errorf("createdAt unmatch"))
	}

	if userMap["objectId"].(string) != user.ID {
		t.Fatal(fmt.Errorf("objectId unmatch"))
	}
}

func TestUserToStruct(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	userStruct := new(CustomUser)

	user.ToStruct(userStruct)

	if userStruct.Username != username {
		t.Fatal(fmt.Errorf("username unmatch"))
	}
}
