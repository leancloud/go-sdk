package leancloud

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

var c *Client

func TestMain(m *testing.M) {
	c = NewEnvClient()

	os.Exit(m.Run())
}

func TestObjectRefCreate(t *testing.T) {
	todo := struct {
		Title      string    `json:"title"`
		Priority   int       `json:"priority"`
		Done       bool      `json:"done"`
		Progress   float64   `json:"progress"`
		FinishedAt time.Time `json:"finishedAt"`
	}{
		Title:      "Team Meeting",
		Priority:   1,
		Done:       false,
		Progress:   12.5,
		FinishedAt: time.Now(),
	}
	/*
		correctResp := map[string]interface{}{
			"title":      "Team Meeting",
			"priority":   1,
			"done":       false,
			"progress":   12.5,
			"finishedAt": todo.FinishedAt,
		}
	*/
	ref, err := c.Class("Todo").Create(todo)
	if err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprint("/1.1/classes/Todo/", ref.ID)
	resp, err := c.request(ServiceAPI, methodGet, path, ref.c.getRequestOptions())
	if err != nil {
		t.Fatal(err)
	}
	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		t.Fatal(err)
	}

	marshalJSON, _ := json.MarshalIndent(respJSON, "", "  ")
	t.Log("\n" + string(marshalJSON))

}

func (object *Object) String() string {
	marshaledJSON, _ := json.MarshalIndent(object, "", "  ")
	return string(marshaledJSON)
}

func TestObjectRefGet(t *testing.T) {
	todo := struct {
		Title      string    `json:"title"`
		Priority   int       `json:"priority"`
		Done       bool      `json:"done"`
		Progress   float64   `json:"progress"`
		FinishedAt time.Time `json:"finishedAt"`
	}{
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

	t.Log(fmt.Sprint("\n", object))
}

func TestObjectRefSet(t *testing.T) {
	todo := struct {
		Title      string    `json:"title"`
		Priority   int       `json:"priority"`
		Done       bool      `json:"done"`
		Progress   float64   `json:"progress"`
		FinishedAt time.Time `json:"finishedAt"`
	}{
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

	testSet := map[string]interface{}{
		"title":      "Another Team Meeting",
		"priority":   10,
		"done":       true,
		"progress":   100.0,
		"finishedAt": time.Now(),
	}

	t.Run("Set String", func(t *testing.T) {
		if err := ref.Set("title", testSet["title"]); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Set Integer", func(t *testing.T) {
		if err := ref.Set("priority", testSet["priority"]); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Set Boolean", func(t *testing.T) {
		if err := ref.Set("done", testSet["done"]); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Set Float", func(t *testing.T) {
		if err := ref.Set("progress", testSet["progress"]); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Set Date", func(t *testing.T) {
		if err := ref.Set("finishedAt", testSet["finishedAt"]); err != nil {
			t.Fatal(err)
		}
	})
}

func TestObjectRefUpdate(t *testing.T) {
	todo := struct {
		Title      string    `json:"title"`
		Priority   int       `json:"priority"`
		Done       bool      `json:"done"`
		Progress   float64   `json:"progress"`
		FinishedAt time.Time `json:"finishedAt"`
	}{
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

	updateMap := map[string]interface{}{
		"title":      "Another Team Meeting",
		"priority":   10,
		"done":       true,
		"progress":   100.0,
		"finishedAt": time.Now(),
	}

	updateStruct := struct {
		Title      string    `json:"title"`
		Priority   int       `json:"priority"`
		Done       bool      `json:"done"`
		Progress   float64   `json:"progress"`
		FinishedAt time.Time `json:"finishedAt"`
	}{
		Title:      "Another Team Meeting",
		Priority:   10,
		Done:       true,
		Progress:   100,
		FinishedAt: time.Now(),
	}

	t.Run("Update with Map", func(t *testing.T) {
		if err := ref.Update(updateMap); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Update with Struct", func(t *testing.T) {
		if err := ref.Update(updateStruct); err != nil {
			t.Fatal(err)
		}
	})
}

func TestObjectRefUpdateWithQuery(t *testing.T) {
	// TODO
}

func TestObjectRefDelete(t *testing.T) {
	todo := struct {
		Title      string    `json:"title"`
		Priority   int       `json:"priority"`
		Done       bool      `json:"done"`
		Progress   float64   `json:"progress"`
		FinishedAt time.Time `json:"finishedAt"`
	}{
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

	if err := ref.Destroy(); err != nil {
		t.Fatal(err)
	}
}
