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
		if client.region == Region(region) {
			if client.appID == appID {
				if client.appKey == appKey {
					if client.masterKey == masterKey {
					} else {
						t.Fatal(errors.New("LEANCLOUD_APP_MASTER_KEY unmatch"))
					}
				} else {
					t.Fatal(errors.New("LEANCLOUD_APP_KEY unmatch"))
				}
			} else {
				t.Fatal(errors.New("LEANCLOUD_APP_ID unmatch"))
			}
		} else {
			t.Fatal(errors.New("LEANCLOUD_REGION unmatch"))
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
		if client.region == Region(region) {
			if client.appID == appID {
				if client.appKey == appKey {
					if client.masterKey == masterKey {
						if client.requestLogger != nil {
						} else {
							t.Fatal(errors.New("unable to set logger"))
						}
					} else {
						t.Fatal(errors.New("LEANCLOUD_APP_MASTER_KEY unmatch"))
					}
				} else {
					t.Fatal(errors.New("LEANCLOUD_APP_KEY unmatch"))
				}
			} else {
				t.Fatal(errors.New("LEANCLOUD_APP_ID unmatch"))
			}
		} else {
			t.Fatal(errors.New("LEANCLOUD_REGION unmatch"))
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
		if client.region == Region(region) {
			if client.appID == appID {
				if client.appKey == appKey {
					if client.masterKey == masterKey {
					} else {
						t.Fatal(errors.New("LEANCLOUD_APP_MASTER_KEY unmatch"))
					}
				} else {
					t.Fatal(errors.New("LEANCLOUD_APP_KEY unmatch"))
				}
			} else {
				t.Fatal(errors.New("LEANCLOUD_APP_ID unmatch"))
			}
		} else {
			t.Fatal(errors.New("LEANCLOUD_REGION unmatch"))
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
		if client.region == Region(region) {
			if client.appID == appID {
				if client.appKey == appKey {
					if client.masterKey == masterKey {
						if client.requestLogger != nil {
						} else {
							t.Fatal(errors.New("unable to set logger"))
						}
					} else {
						t.Fatal(errors.New("LEANCLOUD_APP_MASTER_KEY unmatch"))
					}
				} else {
					t.Fatal(errors.New("LEANCLOUD_APP_KEY unmatch"))
				}
			} else {
				t.Fatal(errors.New("LEANCLOUD_APP_ID unmatch"))
			}
		} else {
			t.Fatal(errors.New("LEANCLOUD_REGION unmatch"))
		}
	})
}

func TestClientClass(t *testing.T) {
	client := &Client{}
	class := client.Class("class")
	if class.c == client {
		if class.Name == "class" {
		} else {
			t.Fatal(errors.New("name of class unmatch"))
		}
	} else {
		t.Fatal(errors.New("client unmatch"))
	}
}

func TestClientObject(t *testing.T) {
	client := &Client{}
	ref := client.Object("class", "f47ac10b58cc4372a5670e02b2c3d479 ")
	if ref.c == client {
		if ref.class == "class" {
			if ref.ID == "f47ac10b58cc4372a5670e02b2c3d479 " {
			} else {
				t.Fatal(errors.New("ID unmatch"))
			}
		} else {
			t.Fatal(errors.New("name of class unmatch"))
		}
	} else {
		t.Fatal(errors.New("client unmatch"))
	}
}
