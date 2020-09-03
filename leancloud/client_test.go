package leancloud

import (
	"errors"
	"os"
	"testing"
)

func TestNewClient(t *testing.T) {
	region, appID, appKey, masterKey := os.Getenv("LEANCLOUD_REGION"), os.Getenv("LEANCLOUD_APP_ID"), os.Getenv("LEANCLOUD_APP_KEY"), os.Getenv("LEANCLOUD_APP_MASTER_KEY")
	t.Run("Production", func(t *testing.T) {
		client := NewClient(region, appID, appKey, masterKey)
		if client == nil {
			t.Fatal(errors.New("unable to create a client"))
		}
		if client.region != Region(region) {
			t.Fatal(errors.New("LEANCLOUD_REGION unmatch"))
		}
		if client.appID != appID {
			t.Fatal(errors.New("LEANCLOUD_APP_ID unmatch"))
		}
		if client.appKey != appKey {
			t.Fatal(errors.New("LEANCLOUD_APP_KEY unmatch"))
		}
		if client.masterKey != masterKey {
			t.Fatal(errors.New("LEANCLOUD_APP_MASTER_KEY unmatch"))
		}
	})

	t.Run("Debug", func(t *testing.T) {
		if err := os.Setenv("LEANCLOUD_DEBUG", "true"); err != nil {
			t.Fatal("unable to set debugging flag")
		}
		client := NewClient(region, appID, appKey, masterKey)
		if client == nil {
			t.Fatal(errors.New("unable to create a client"))
		}
		if client.region != Region(region) {
			t.Fatal(errors.New("LEANCLOUD_REGION unmatch"))
		}
		if client.appID != appID {
			t.Fatal(errors.New("LEANCLOUD_APP_ID unmatch"))
		}
		if client.appKey != appKey {
			t.Fatal(errors.New("LEANCLOUD_APP_KEY unmatch"))
		}
		if client.masterKey != masterKey {
			t.Fatal(errors.New("LEANCLOUD_APP_MASTER_KEY unmatch"))
		}
		if client.requestLogger == nil {
			t.Fatal(errors.New("unable to set logger"))
		}
	})
}

func TestNewEnvClient(t *testing.T) {
	region, appID, appKey, masterKey := os.Getenv("LEANCLOUD_REGION"), os.Getenv("LEANCLOUD_APP_ID"), os.Getenv("LEANCLOUD_APP_KEY"), os.Getenv("LEANCLOUD_APP_MASTER_KEY")
	t.Run("Production", func(t *testing.T) {
		client := NewEnvClient()
		if client == nil {
			t.Fatal(errors.New("unable to create a client"))
		}
		if client.region != Region(region) {
			t.Fatal(errors.New("LEANCLOUD_REGION unmatch"))
		}
		if client.appID != appID {
			t.Fatal(errors.New("LEANCLOUD_APP_ID unmatch"))
		}
		if client.appKey != appKey {
			t.Fatal(errors.New("LEANCLOUD_APP_KEY unmatch"))
		}
		if client.masterKey != masterKey {
			t.Fatal(errors.New("LEANCLOUD_APP_MASTER_KEY unmatch"))
		}
	})

	t.Run("Debug", func(t *testing.T) {
		if err := os.Setenv("LEANCLOUD_DEBUG", "true"); err != nil {
			t.Fatal("unable to set debugging flag")
		}
		client := NewEnvClient()
		if client == nil {
			t.Fatal(errors.New("unable to create a client"))
		}
		if client.region != Region(region) {
			t.Fatal(errors.New("LEANCLOUD_REGION unmatch"))
		}
		if client.appID != appID {
			t.Fatal(errors.New("LEANCLOUD_APP_ID unmatch"))
		}
		if client.appKey != appKey {
			t.Fatal(errors.New("LEANCLOUD_APP_KEY unmatch"))
		}
		if client.masterKey != masterKey {
			t.Fatal(errors.New("LEANCLOUD_APP_MASTER_KEY unmatch"))
		}
		if client.requestLogger == nil {
			t.Fatal(errors.New("unable to set logger"))
		}
	})
}

func TestClientClass(t *testing.T) {
	client := &Client{}
	class := client.Class("class")
	if class.c != client {
		t.Fatal(errors.New("client unmatch"))
	}
	if class.Name != "class" {
		t.Fatal(errors.New("name of class unmatch"))
	}
}

func TestClientObject(t *testing.T) {
	client := &Client{}
	ref := client.Object("class", "f47ac10b58cc4372a5670e02b2c3d479")
	if ref.c != client {
		t.Fatal(errors.New("client unmatch"))
	}
	if ref.class != "class" {
		t.Fatal(errors.New("name of class unmatch"))
	}
	if ref.ID != "f47ac10b58cc4372a5670e02b2c3d479" {
		t.Fatal(errors.New("ID unmatch"))
	}
}
