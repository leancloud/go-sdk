package leancloud

import (
	"fmt"
	"math/rand"
	"testing"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomStringBytes(n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	return string(b)
}

func TestUsersSignUp(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	if user.fields["username"].(string) != username {
		t.Fatal(fmt.Errorf("username unmatch"))
	}
}

func TestUsersLogIn(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	loginUser, err := c.Users.LogIn(username, password)
	if err != nil {
		t.Fatal(err)
	}

	if user.fields["username"].(string) != loginUser.fields["username"].(string) {
		t.Fatal(fmt.Errorf("username unmatch"))
	}

	if user.ID != loginUser.ID {
		t.Fatal(fmt.Errorf("objectId unmatch"))
	}
}

func TestUsersBecome(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	sessionUser, err := c.Users.Become(user.sessionToken)
	if err != nil {
		t.Fatal(err)
	}

	if user.ID != sessionUser.ID {
		t.Fatal(fmt.Errorf("objectId unmatch"))
	}

	if user.sessionToken != sessionUser.sessionToken {
		t.Fatal(fmt.Errorf("sessionToken unmatch"))
	}
}
