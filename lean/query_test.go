package lean

import (
	"testing"
)

func Test_Query_CloudQuery(t *testing.T) {
	users := []User{}

	err := client.CloudQuery("select * from _User limit ?", []interface{}{2}, &users, UseMasterKey(true))

	if err != nil {
		t.Error(err)
	}

	t.Log(users)

	if len(users) > 2 {
		t.Error("len(users) > 2")
	}
}

func Test_Query_Get(t *testing.T) {
	user := User{
		ObjectMeta: ObjectID("56ef63d07db2a20052226804"),
	}

	err := client.NewQuery("_User").Get(&user)

	if err != nil {
		t.Error(err)
	}

	t.Log(user)
}

func Test_Query_Find(t *testing.T) {
	users := []User{}

	err := client.NewQuery("_User").Find(&users, UseMasterKey(true))

	if err != nil {
		t.Error(err)
	}

	t.Log(users)
}

func Test_Query_First(t *testing.T) {
	user := User{}

	err := client.NewQuery("_User").First(&user, UseMasterKey(true))

	if err != nil {
		t.Error(err)
	}

	t.Log(user)
}
