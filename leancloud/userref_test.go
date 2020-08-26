package leancloud

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestUserRefGet(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	userByID, err := c.User(user.ID).Get()
	if err != nil {
		t.Fatal(err)
	}

	if userByID.ID != user.ID {
		t.Fatal(fmt.Errorf("objectId unmatch"))
	}

	if user.fields["username"].(string) != userByID.fields["username"].(string) {
		t.Fatal(fmt.Errorf("username unmatch"))
	}
}

func TestUserRefSet(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.User(user.ID).Set("set", 1, UseSessionToken(user.sessionToken)); err != nil {
		t.Fatal(err)
	}

	userByID, err := c.User(user.ID).Get()
	if err != nil {
		t.Fatal(err)
	}

	if int(userByID.fields["set"].(float64)) != 1 {
		t.Fatal("set unmatch")
	}
}

func TestUserRefUpdate(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	if err := c.User(user.ID).Update(map[string]interface{}{
		"content":     "extra string",
		"number":      10,
		"done":        true,
		"floatNumber": 100.5,
		"extraDate":   now,
	}, UseSessionToken(user.sessionToken)); err != nil {
		t.Fatal(err)
	}

	userByID, err := c.User(user.ID).Get()
	if userByID.fields["content"].(string) != "extra string" {
		t.Fatal(fmt.Errorf("content unmatch"))
	}

	if int(userByID.fields["number"].(float64)) != 10 {
		t.Fatal(fmt.Errorf("number unmatch"))
	}

	if userByID.fields["done"].(bool) != true {
		t.Fatal(fmt.Errorf("done unmatch"))
	}

	if userByID.fields["floatNumber"].(float64) != 100.5 {
		t.Fatal(fmt.Errorf("floatNumber unmatch"))
	}

	extraDate := userByID.fields["extraDate"].(time.Time)

	if extraDate.Unix() != now.Unix() {
		t.Fatal(fmt.Errorf("extraDate unmatch"))
	}
}

func TestUserRefUpdateWithQuery(t *testing.T) {
	// TODO
}

func TestUserRefDestroy(t *testing.T) {
	username, password := randomStringBytes(10), randomStringBytes(20)
	user, err := c.Users.SignUp(username, password)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.User(user.ID).Destroy(UseMasterKey(true)); err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprint("/1.1/users/", user.ID)
	_, err = c.request(ServiceAPI, methodGet, path, c.getRequestOptions())
	if !strings.Contains(err.Error(), "Could not find user.") {
		t.Fatal(err)
	}
}
