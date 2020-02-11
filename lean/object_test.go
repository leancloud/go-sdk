package lean

import (
	"testing"
	"time"
)

var client = NewEnvClient()

type Todo struct {
	ObjectMeta

	Title      String  `json:"title"`
	Priority   Integer `json:"priority"`
	Done       Boolean `json:"done"`
	Progress   Float   `json:"progress"`
	FinishedAt Date    `json:"finishedAt"`
}

func (*Todo) ClassName() string {
	return "Todo"
}

func Test_Object_Create(t *testing.T) {
	todo := Todo{
		Title:      NewString("Team meeting"),
		Priority:   NewInteger(5),
		Done:       NewBoolean(true),
		Progress:   NewFloat(12.5),
		FinishedAt: NewDate(time.Now()),
	}

	err := client.Save(&todo)

	if err != nil {
		t.Error(err)
	}

	t.Log(todo)
}

func Test_Object_Destroy(t *testing.T) {
	todo := Todo{
		Title: NewString("Delete this record"),
	}

	err := client.Save(&todo)

	if err != nil {
		t.Error(err)
	}

	err = client.Destroy(&todo)

	if err != nil {
		t.Error(err)
	}
}
