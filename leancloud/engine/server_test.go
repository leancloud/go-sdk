package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/leancloud/go-sdk/leancloud"
	"github.com/levigross/grequests"
)

var cloudEndpoint = "http://localhost:3000"

func TestMain(m *testing.M) {
	go http.ListenAndServe(":3000", Handler(nil))

	Define("hello", func(r *Request) (interface{}, error) {
		return map[string]string{
			"hello": "world",
		}, nil
	})

	DefineWithOption("hello_with_option_internal", func(r *Request) (interface{}, error) {
		return map[string]string{
			"hello": "world",
		}, nil
	}, &DefineOption{
		NotFetchUser: true,
		Internal:     true,
	})

	DefineWithOption("hello_with_option_fetch_user", func(r *Request) (interface{}, error) {
		return map[string]string{
			"sessionToken": r.SessionToken,
		}, nil
	}, &DefineOption{
		NotFetchUser: false,
	})

	DefineWithOption("hello_with_option_not_fetch_user", func(r *Request) (interface{}, error) {
		return map[string]interface{}{
			"currentUser": r.CurrentUser,
		}, nil
	}, &DefineOption{
		NotFetchUser: true,
		Internal:     false,
	})

	os.Exit(m.Run())
}

func TestMetadataResponse(t *testing.T) {
	resp, err := grequests.Get(cloudEndpoint+"/1.1/functions/_ops/metadata", &grequests.RequestOptions{})
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Close()

	metadata := new(metadataResponse)
	if err := json.NewDecoder(resp.RawResponse.Body).Decode(metadata); err != nil {
		t.Fatal(err)
	}

	for _, v := range metadata.Result {
		if functions[v] == nil {
			t.Fatal(fmt.Errorf("cannot found cloud function"))
		}
	}
}

func TestHandler(t *testing.T) {
	t.Run("function call", func(t *testing.T) {
		resp, err := grequests.Get(cloudEndpoint+"/1.1/functions/hello", nil)
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

		if respBody["hello"] != "world" {
			t.Fatal("unexpected response format")
		}
	})

	t.Run("function call with sessionToken", func(t *testing.T) {
		user, err := client.User("5f86a88f27075b72775de082").Get(leancloud.UseMasterKey(true))
		if err != nil {
			t.Fatal(err)
		}

		options := grequests.RequestOptions{
			Headers: map[string]string{
				"X-LC-Id":      os.Getenv("LEANCLOUD_APP_ID"),
				"X-LC-Key":     os.Getenv("LEANCLOUD_APP_KEY"),
				"X-LC-Session": user.GetSessionToken(),
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

		if respBody["sessionToken"] != user.GetSessionToken() {
			t.Fatal("unexpected response format")
		}
	})

	t.Run("function call shoud not found", func(t *testing.T) {
		resp, err := grequests.Get(cloudEndpoint+"/1.1/functions/not_found", nil)
		if err != nil {
			if resp.StatusCode != http.StatusNotFound {
				t.Fatal(err)
			}
		}
	})
}
