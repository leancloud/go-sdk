package leancloud

import (
	"fmt"
	"testing"
	"time"
)

func TestObjectGetMap(t *testing.T) {
	todo := Todo{
		Title:      "Team Meeting",
		Priority:   1,
		Done:       false,
		Progress:   12.5,
		FinishedAt: time.Now(),
	}

	ref, err := c.Class("Todo").Create(todo)
	if err != nil {
		t.Fatal(err)
	}

	object, err := ref.Get()
	if err != nil {
		t.Fatal(err)
	}

	mapObject := object.GetMap()

	if mapObject["title"].(string) != todo.Title {
		t.Fatal(fmt.Errorf("title unmatch"))
	}

	if (int)(mapObject["priority"].(float64)) != todo.Priority {
		t.Fatal(fmt.Errorf("priority unmatch"))
	}

	if mapObject["done"].(bool) != todo.Done {
		t.Fatal(fmt.Errorf("done unmatch"))
	}

	if mapObject["progress"].(float64) != todo.Progress {
		t.Fatal(fmt.Errorf("progress unmatch"))
	}

	finishedAt, ok := mapObject["finishedAt"].(time.Time)
	if !ok {
		t.Fatal("unable to parse finishedAt from mapObject")
	}
	if finishedAt.Unix() != todo.FinishedAt.Unix() {
		t.Fatal(fmt.Errorf("finishedAt unmatch"))
	}

	if mapObject["objectId"].(string) != ref.ID {
		t.Fatal("objectId unmatch")
	}

	t.Log(mapObject)
}

func TestObjectToStruct(t *testing.T) {
	todo := Todo{
		Title:      "Team Meeting",
		Priority:   1,
		Done:       false,
		Progress:   12.5,
		FinishedAt: time.Now(),
	}

	ref, err := c.Class("Todo").Create(todo)
	if err != nil {
		t.Fatal(err)
	}

	object, err := ref.Get()
	if err != nil {
		t.Fatal(err)
	}

	dup := new(Todo)
	object.ToStruct(dup)

	if dup.Title != todo.Title {
		t.Fatal(fmt.Errorf("title unmatch"))
	}

	if dup.Priority != todo.Priority {
		t.Fatal(fmt.Errorf("priority unmatch"))
	}

	if dup.Done != todo.Done {
		t.Fatal(fmt.Errorf("done unmatch"))
	}

	if dup.Progress != todo.Progress {
		t.Fatal(fmt.Errorf("unable to parse finishedAt from mapObject"))
	}

	if dup.FinishedAt.Unix() != todo.FinishedAt.Unix() {
		t.Fatal(fmt.Errorf("finishedAt unmatch"))
	}

	t.Log(dup)
}
