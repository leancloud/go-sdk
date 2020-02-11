package lean

import (
	"os"
	"testing"
)

var alice User

func TestMain(m *testing.M) {
	code := m.Run()

	if err := client.Destroy(&alice, UseMasterKey(true)); err != nil {
		panic(err)
	}

	os.Exit(code)
}

func Test_User_Save(t *testing.T) {
	alice = User{
		Username: NewString("alice"),
		Password: NewString("blabla"),
	}

	err := client.Save(&alice)

	if err != nil {
		t.Error(err)
	}

	t.Log(alice)
}
