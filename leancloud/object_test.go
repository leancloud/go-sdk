package leancloud

import (
	"fmt"
	"testing"
	"time"
)

func TestObjectGetMap(t *testing.T) {
	todo := Todo{
		Title:        "Team Meeting",
		Priority:     1,
		Done:         false,
		Progress:     12.5,
		FinishedAt:   time.Now(),
		Participants: []string{"Adams", "Baker", "Clark", "Davis", "Evans", "Frank"},
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
}

type embedObject struct {
	Name     string  `json:"name"`
	Age      int     `json:"age"`
	Progress float32 `json:"progess"`
}

func TestObjectToStruct(t *testing.T) {
	todo := Todo{
		Title:        "Team Meeting",
		Priority:     1,
		Done:         false,
		Progress:     12.5,
		FinishedAt:   time.Now(),
		Participants: []string{"Adams", "Baker", "Clark", "Davis", "Evans", "Frank"},
		Dates:        []time.Time{time.Now(), time.Now(), time.Now()},
		Objects: []embedObject{
			{
				Name:     "1",
				Age:      11,
				Progress: 11.5,
			},
			{
				Name:     "2",
				Age:      12,
				Progress: 12.5,
			},
		},
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
	if err := object.ToStruct(dup); err != nil {
		t.Fatal(err)
	}

}
