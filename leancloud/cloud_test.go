package leancloud

import (
	"fmt"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	resp, err := Run("hello", nil)
	if err != nil {
		t.Fatal(err)
	}

	respString, ok := resp.(string)
	if !ok {
		t.Fatal("unexpected response format")
	}

	if respString != "Hello world!" {
		t.Fatal("unexpected response format")
	}
}

func TestRunWithOptions(t *testing.T) {
	t.Run("local", func(t *testing.T) {
		resp, err := RunWithOption("hello", nil, &RunOption{
			Remote:       false,
			User:         nil,
			SessionToken: "",
		})
		if err != nil {
			t.Fatal(err)
		}

		respMap, ok := resp.(map[string]string)
		if !ok {
			t.Fatal(fmt.Errorf("unmatch response"))
		}

		if respMap["hello"] != "world" {
			t.Fatal(fmt.Errorf("unmatch response"))
		}

	})

	t.Run("hello_with_option_internal", func(t *testing.T) {
		t.Run("remote", func(t *testing.T) {
			_, err := RunWithOption("hello_with_option_internal", nil, &RunOption{
				Remote: true,
			})

			if err != nil {
				if !strings.Contains(err.Error(), "401 Internal cloud function") {
					t.Fatal(err)
				}
			}
		})

		t.Run("local", func(t *testing.T) {
			resp, err := RunWithOption("hello_with_option_internal", nil, &RunOption{
				Remote: false,
			})

			if err != nil {
				t.Fatal(err)
			}

			respMap, ok := resp.(map[string]string)
			if !ok {
				t.Fatal(fmt.Errorf("unmatch response"))
			}

			if respMap["hello"] != "world" {
				t.Fatal(fmt.Errorf("unmatch response"))
			}
		})
	})

	t.Run("hello_with_option_fetch_user", func(t *testing.T) {
		/*
			user, err := client.User("5fa504d0f98fd535ebe8b3f0").Get(UseMasterKey(true))
			if err != nil {
				t.Fatal(err)
			}

			t.Run("remote", func(t *testing.T) {
				resp, err := RunWithOption("hello_with_option_fetch_user", nil, &RunOption{
					Remote:       true,
					SessionToken: user.SessionToken(),
				})

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

				if sessionToken != user.SessionToken() {
					t.Fatal("unexpected response format")
				}

			})

			t.Run("local", func(t *testing.T) {
				resp, err := RunWithOption("hello_with_option_fetch_user", nil, &RunOption{
					Remote:       false,
					SessionToken: user.SessionToken(),
				})

				if err != nil {
					t.Fatal(err)
				}

				respMap, ok := resp.(map[string]string)
				if !ok {
					t.Fatal("unexpected response format")
				}

				if respMap["sessionToken"] != user.SessionToken() {
					t.Fatal("unexpected response format")
				}
			})
		*/
	})

	t.Run("don't fetch user", func(t *testing.T) {
		t.Run("remote", func(t *testing.T) {
			resp, err := RunWithOption("hello_with_option_not_fetch_user", nil, &RunOption{
				Remote: true,
			})
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
			resp, err := RunWithOption("hello_with_option_not_fetch_user", nil, &RunOption{
				Remote: true,
			})
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
			_, err := RunWithOption("not_found", nil, &RunOption{
				Remote: true,
			})

			if err != nil {
				if !strings.Contains(err.Error(), "No such cloud function") {
					t.Fatal(err)
				}
			}
		})

		t.Run("local", func(t *testing.T) {
			_, err := RunWithOption("not_found", nil, &RunOption{
				Remote: false,
			})

			if err != nil {
				if !strings.Contains(err.Error(), "no such cloud function") {
					t.Fatal(err)
				}
			}
		})
	})
}
