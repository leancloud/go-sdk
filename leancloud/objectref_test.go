package leancloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
)

type Todo struct {
	Title      string    `json:"title"`
	Priority   int       `json:"priority"`
	Done       bool      `json:"done"`
	Progress   float64   `json:"progress"`
	FinishedAt time.Time `json:"finishedAt"`
}

var c *Client

func TestMain(m *testing.M) {
	c = NewEnvClient()
	rand.Seed(time.Now().UnixNano())

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

func TestObjectRefCreate(t *testing.T) {
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

	path := fmt.Sprint("/1.1/classes/Todo/", ref.ID)
	resp, err := c.request(ServiceAPI, MethodGet, path, ref.c.getRequestOptions())
	if err != nil {
		t.Fatal(err)
	}
	respJSON := map[string]interface{}{}
	if err := json.Unmarshal(resp.Bytes(), &respJSON); err != nil {
		t.Fatal(err)
	}
	if respJSON["title"].(string) != todo.Title {
		t.Fatal(errors.New("value of title unmatch"))
	}
	if (int)(respJSON["priority"].(float64)) != todo.Priority {
		t.Fatal(errors.New("value of priority unmatch"))
	}
	if respJSON["done"].(bool) != todo.Done {
		t.Fatal(errors.New("value of done unmatch"))
	}
	finishedAt, _ := respJSON["finishedAt"].(map[string]interface{})
	date, _ := decodeDate(finishedAt)
	if date.Unix() != todo.FinishedAt.Unix() {
		t.Fatal(errors.New("value of finishedAt field unmatch"))
	}
}

func TestObjectRefGet(t *testing.T) {
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

	if object.fields["title"].(string) != todo.Title {
		t.Fatal(errors.New("value of title unmatch"))
	}
	if (int)(object.fields["priority"].(float64)) != todo.Priority {
		t.Fatal(errors.New("value of priority unmatch"))
	}
	if object.fields["done"].(bool) != todo.Done {
		t.Fatal(errors.New("value of done field unmatch"))
	}
	if object.fields["progress"].(float64) != todo.Progress {
		t.Fatal(errors.New("value of progress field unmatch"))
	}
	finishedAt := object.fields["finishedAt"].(time.Time)
	if finishedAt.Unix() != todo.FinishedAt.Unix() {
		t.Fatal(errors.New("value of finishedAt field unmatch"))
	}
}

func TestObjectRefSet(t *testing.T) {
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

	if err := ref.Set("title", "Another Team Meeting"); err != nil {
		t.Fatal(err)
	}

	object, err := c.Class("Todo").Object(ref.ID).Get()
	if err != nil {
		t.Fatal(err)
	}

	if object.fields["title"].(string) != "Another Team Meeting" {
		t.Fatal("value of title unchanged")
	}
}

func TestObjectRefUpdate(t *testing.T) {
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

	updateMap := map[string]interface{}{
		"title":      "Another Team Meeting",
		"priority":   10,
		"done":       true,
		"progress":   100.0,
		"finishedAt": time.Now(),
	}

	if err := ref.Update(updateMap); err != nil {
		t.Fatal(err)
	}

	object, err := c.Class("Todo").Object(ref.ID).Get()
	if err != nil {
		t.Fatal(err)
	}
	if object.fields["title"].(string) != updateMap["title"] {
		t.Fatal(errors.New("value of title unmatch"))
	}
	if (int)(object.fields["priority"].(float64)) != updateMap["priority"] {
		t.Fatal(errors.New("value of priority unmatch"))
	}
	if object.fields["done"].(bool) != updateMap["done"] {
		t.Fatal(errors.New("value of done field unmatch"))
	}
	if object.fields["progress"].(float64) != updateMap["progress"] {
		t.Fatal(errors.New("value of progress field unmatch"))
	}
	finishedAt := object.fields["finishedAt"].(time.Time)
	if finishedAt.Unix() != updateMap["finishedAt"].(time.Time).Unix() {
		t.Fatal(errors.New("value of finishedAt field unmatch"))
	}
}

func TestObjectRefUpdateWithQuery(t *testing.T) {
	// TODO
}

func TestObjectRefDestroy(t *testing.T) {
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

	if err := ref.Destroy(); err != nil {
		t.Fatal(err)
	}

	path := fmt.Sprint("/1.1/classes/Todo/", ref.ID)
	resp, err := c.request(ServiceAPI, MethodGet, path, c.getRequestOptions())
	if err != nil {
		t.Fatal(err)
	}
	if string(resp.Bytes()) != "{}" {
		t.Fatal("unable to destroy object")
	}
}
