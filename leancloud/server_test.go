package leancloud

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/levigross/grequests"
)

var cloudEndpoint = "http://localhost:3000"

func TestMain(m *testing.M) {
	go http.ListenAndServe(":3000", Engine.Handler())

	os.Exit(m.Run())
}
func TestMetadataResponse(t *testing.T) {
	resp, err := grequests.Get(cloudEndpoint+"/1.1/functions/_ops/metadatas", &grequests.RequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	metadata := new(metadataResponse)
	if err := json.NewDecoder(resp.RawResponse.Body).Decode(metadata); err != nil {
		if err != io.EOF {
			t.Fatal(err)
		}
	}

	for _, v := range metadata.Result {
		if Engine.(*engine).functions[v] == nil {
			t.Fatal(fmt.Errorf("cannot found cloud function"))
		}
	}
}

func TestHandler(t *testing.T) {
	t.Run("function call", func(t *testing.T) {
		resp, err := grequests.Get(cloudEndpoint+"/1.1/functions/hello", &grequests.RequestOptions{
			Headers: map[string]string{
				"X-LC-Id":  os.Getenv("LEANCLOUD_APP_ID"),
				"X-LC-Key": os.Getenv("LEANCLOUD_APP_KEY"),
			},
		})
		if err != nil {
			t.Fatal(err)
		}

		ret := new(functionResponse)
		if err := json.Unmarshal(resp.Bytes(), ret); err != nil {
			t.Log(string(resp.Bytes()))
			t.Fatal(err)
		}

		respBody, ok := ret.Result.(map[string]interface{})
		if !ok {
			t.Fatal("unexpected response format")
		}

		if respBody["Hello"] != "World" {
			t.Fatal("unexpected response format")
		}
	})

	t.Run("function call with sessionToken", func(t *testing.T) {
		user := new(User)
		if err := client.Users.ID(testUserID).Get(user, UseMasterKey(true)); err != nil {
			t.Fatal(err)
		}

		options := grequests.RequestOptions{
			Headers: map[string]string{
				"X-LC-Id":      os.Getenv("LEANCLOUD_APP_ID"),
				"X-LC-Key":     os.Getenv("LEANCLOUD_APP_KEY"),
				"X-LC-Session": user.SessionToken,
			},
		}

		resp, err := grequests.Get(cloudEndpoint+"/1.1/functions/hello_with_option_fetch_user", &options)
		if err != nil {
			t.Fatal(err)
		}

		ret := new(functionResponse)
		if err := json.Unmarshal(resp.Bytes(), ret); err != nil {
			t.Log(string(resp.Bytes()))
			t.Fatal(err)
		}

		respBody, ok := ret.Result.(map[string]interface{})
		if !ok {
			t.Fatal("unexpected response format")
		}

		if respBody["sessionToken"] != user.SessionToken {
			t.Fatal("unexpected response format")
		}
	})

	t.Run("function call should not found", func(t *testing.T) {
		resp, err := grequests.Get(cloudEndpoint+"/1.1/functions/not_found", nil)
		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				t.Fatal(err)
			}
		}
	})
}
