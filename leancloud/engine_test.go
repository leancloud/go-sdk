package leancloud

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

var testUserID string

func init() {
	testUserID = os.Getenv("TEST_USER_ID")
	Engine.Define("hello", func(r *FunctionRequest) (interface{}, error) {
		return map[string]string{
			"Hello": "World",
		}, nil
	})

	Engine.Define("hello_with_option_internal", func(r *FunctionRequest) (interface{}, error) {
		return map[string]string{
			"Hello": "World",
		}, nil
	}, WithInternal(), WithoutFetchUser())

	Engine.Define("hello_with_option_fetch_user", func(r *FunctionRequest) (interface{}, error) {
		return map[string]string{
			"sessionToken": r.SessionToken,
		}, nil
	})

	Engine.Define("hello_with_option_not_fetch_user", func(r *FunctionRequest) (interface{}, error) {
		return map[string]interface{}{
			"sessionToken": r.SessionToken,
		}, nil
	}, WithoutFetchUser())

	Engine.Define("hello_with_object", func(r *FunctionRequest) (interface{}, error) {
		return r.CurrentUser, nil
	})
}

func TestRun(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		resp, err := Engine.Run("hello", nil)
		if err != nil {
			t.Fatal(err)
		}

		respMap, ok := resp.(map[string]string)
		if !ok {
			t.Fatal(fmt.Errorf("unmatch response"))
		}

		if respMap["Hello"] != "World" {
			t.Fatal(fmt.Errorf("unmatch response"))
		}

	})

	t.Run("hello_with_option_internal", func(t *testing.T) {
		t.Run("remote", func(t *testing.T) {
			_, err := client.Run("hello_with_option_internal", nil)

			if err != nil {
				if !strings.Contains(err.Error(), "401 Internal cloud function") {
					t.Fatal(err)
				}
			}
		})

		t.Run("local", func(t *testing.T) {
			resp, err := Engine.Run("hello_with_option_internal", nil)

			if err != nil {
				t.Fatal(err)
			}

			respMap, ok := resp.(map[string]string)
			if !ok {
				t.Fatal(fmt.Errorf("unmatch response"))
			}

			if respMap["Hello"] != "World" {
				t.Fatal(fmt.Errorf("unmatch response"))
			}
		})
	})

	t.Run("hello_with_option_fetch_user", func(t *testing.T) {
		user := new(User)
		if err := client.Users.ID(testUserID).Get(user, UseMasterKey(true)); err != nil {
			t.Fatal(err)
		}

		t.Run("remote", func(t *testing.T) {
			resp, err := client.Run("hello_with_option_fetch_user", nil, WithSessionToken(user.SessionToken))

			if err != nil {
				t.Fatal(err)
			}

			respMap, ok := resp.(map[string]interface{})
			if !ok {
				t.Fatal("unexpected response format")
			}

			sessionToken, ok := respMap["sessionToken"].(string)

			if !ok {
				t.Fatal("unexpected response format")
			}

			if sessionToken != user.SessionToken {
				t.Fatal("unexpected response format")
			}

		})

		t.Run("local", func(t *testing.T) {
			resp, err := Engine.Run("hello_with_option_fetch_user", nil, WithSessionToken(user.SessionToken))

			if err != nil {
				t.Fatal(err)
			}

			respMap, ok := resp.(map[string]string)
			if !ok {
				t.Fatal("unexpected response format")
			}

			if respMap["sessionToken"] != user.SessionToken {
				t.Fatal("unexpected response format")
			}
		})
	})

	t.Run("don't fetch user", func(t *testing.T) {
		t.Run("remote", func(t *testing.T) {
			resp, err := client.Run("hello_with_option_not_fetch_user", nil)
			if err != nil {
				t.Fatal(err)
			}

			respMap, ok := resp.(map[string]interface{})
			if !ok {
				t.Fatal("unexpected response format")
			}

			if len(respMap) != 0 {
				t.Fatal("unexpected response format")
			}
		})

		t.Run("local", func(t *testing.T) {
			resp, err := Engine.Run("hello_with_option_not_fetch_user", nil)
			if err != nil {
				t.Fatal(err)
			}

			respMap, ok := resp.(map[string]interface{})
			if !ok {
				t.Fatal("unexpected response format")
			}

			if respMap["currentUser"] != nil {
				t.Fatal("unexpected response format")
			}
		})
	})

	t.Run("not_found", func(t *testing.T) {
		t.Run("remote", func(t *testing.T) {
			_, err := client.Run("not_found", nil)

			if err != nil {
				if !strings.Contains(err.Error(), "No such cloud function") {
					t.Fatal(err)
				}
			}
		})

		t.Run("local", func(t *testing.T) {
			_, err := Engine.Run("not_found", nil)

			if err != nil {
				if !strings.Contains(err.Error(), "no such cloud function") {
					t.Fatal(err)
				}
			}
		})
	})
}

func TestRPC(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		user := new(User)
		if err := client.Users.ID(testUserID).Get(user, UseMasterKey(true)); err != nil {
			t.Fatal(err)
		}

		retUser := new(User)
		err := Engine.RPC("hello_with_object", nil, retUser, WithUser(user))
		if err != nil {
			t.Fatal(err)
		}

		if retUser.SessionToken != user.SessionToken {
			t.Fatal(fmt.Errorf("dismatch sessionToken"))
		}
	})

	t.Run("remote", func(t *testing.T) {
		user := new(User)
		if err := client.Users.ID(testUserID).Get(user, UseMasterKey(true)); err != nil {
			t.Fatal(err)
		}

		retUser := new(User)
		err := client.RPC("hello_with_object", nil, retUser, WithUser(user))
		if err != nil {
			t.Fatal(err)
		}

		if retUser.ID != user.ID {
			t.Fatal(fmt.Errorf("dismatch sessionToken"))
		}
	})
}
