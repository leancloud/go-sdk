package leancloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

type Todo struct {
	ID           string        `json:"objectId"`
	Title        string        `json:"title"`
	Priority     int           `json:"priority"`
	Done         bool          `json:"done"`
	Progress     float64       `json:"progress"`
	FinishedAt   time.Time     `json:"finishedAt"`
	Participants []string      `json:"participants"`
	Dates        []time.Time   `json:"dates"`
	Objects      []embedObject `json:"objects"`
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
	date, _ := decodeDate(finishedAt["iso"].(string))
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

func TestObjectRefNested(t *testing.T) {
	t.Run("Array", func(t *testing.T) {
		nameList := []string{"Adams", "Baker", "Clark", "Davis", "Evans", "Frank"}
		posts := map[string]interface{}{
			"title": "Hello World!",
			"reply": nameList,
		}

		ref, err := client.Class("Posts").Create(posts)
		if err != nil {
			t.Fatal(err)
		}

		object, err := ref.Get()
		if err != nil {
			t.Fatal(err)
		}

		reply := object.Get("reply")
		infReply, ok := reply.([]interface{})
		if !ok {
			t.Fatal(fmt.Errorf("unexpcted type for infReply: want []interface{} but %v", reflect.TypeOf(reply)))
		}

		replyArray := make([]string, len(infReply))
		for i, d := range infReply {
			replyArray[i] = d.(string)
		}

		if !reflect.DeepEqual(replyArray, nameList) {
			t.Fatal(fmt.Errorf("response is not match with the original content"))
		}
	})

	t.Run("Object", func(t *testing.T) {
		adams := map[string]interface{}{
			"name":   "Adams",
			"gender": "Male",
			"age":    22,
		}

		ref, err := client.Class("People").Create(adams)
		if err != nil {
			t.Fatal(err)
		}

		adamsObject, err := ref.Get()
		if err != nil {
			t.Fatal(err)
		}

		adamsObject.isPointer = true
		adamsObject.fields["className"] = "People"

		nameList := []string{"Adams", "Baker", "Clark", "Davis", "Evans", "Frank"}
		posts := map[string]interface{}{
			"title":  "Hello World!",
			"author": adamsObject,
			"reply":  nameList,
		}

		postsRef, err := client.Class("Posts").Create(posts)
		if err != nil {
			t.Fatal(err)
		}

		postsObject, err := postsRef.Get()
		if err != nil {
			t.Fatal(err)
		}

		if reflect.TypeOf(postsObject.fields["author"]) != reflect.TypeOf(&Object{}) {
			t.Fatal(fmt.Errorf("type of field of author shoud bd *Object but %v", reflect.TypeOf(postsObject.fields["author"])))
		}
	})

	t.Run("GeoPoint", func(t *testing.T) {
		cities := map[string]interface{}{
			"location": &GeoPoint{
				Latitude:  31.385597,
				Longitude: 120.980736,
			},
		}

		ref, err := client.Class("Cities").Create(cities)
		if err != nil {
			t.Fatal(err)
		}

		object, err := ref.Get()
		if err != nil {
			t.Fatal(err)
		}

		if reflect.TypeOf(object.fields["location"]) != reflect.TypeOf(&GeoPoint{}) {
			t.Fatal(fmt.Errorf("type of field of location shoud bd *GeoPoint but %v", reflect.TypeOf(object.fields["location"])))
		}
	})

	t.Run("Byte", func(t *testing.T) {
		data := map[string]interface{}{
			"data": []byte("Hello Type Byte!"),
		}

		ref, err := client.Class("Data").Create(data)
		if err != nil {
			t.Fatal(err)
		}

		object, err := ref.Get()
		if err != nil {
			t.Fatal(err)
		}

		if reflect.TypeOf(object.fields["data"]) != reflect.TypeOf([]byte{}) {
			t.Fatal(fmt.Errorf("type of field of location shoud bd []byte but %v", reflect.TypeOf(object.fields["data"])))
		}
	})

	t.Run("Date", func(t *testing.T) {
		date := map[string]interface{}{
			"current": time.Now(),
		}

		ref, err := client.Class("Date").Create(date)
		if err != nil {
			t.Fatal(err)
		}

		object, err := ref.Get()

		if reflect.TypeOf(object.fields["current"]) != reflect.TypeOf(time.Time{}) {
			t.Fatal(fmt.Errorf("type of field of current shoud bd []byte but %v", reflect.TypeOf(object.fields["current"])))
		}
	})

	t.Run("File", func(t *testing.T) {
		filename, err := generateTempFile("go-sdk-file-upload-*.txt")
		if err != nil {
			t.Fatal(err)
		}

		defer os.Remove(filename)

		fd, err := os.Open(filename)
		if err != nil {
			t.Fatal(err)
		}

		_, name := filepath.Split(filename)
		file := &File{
			Name: name,
			MIME: mime.TypeByExtension(filepath.Ext(name)),
		}

		if err := c.Files.Upload(file, fd); err != nil {
			t.Fatal(err)
		}

		if err := checkUpload(file.URL); err != nil {
			t.Fatal(err)
		}

		fileMap := map[string]interface{}{
			"file": file,
		}

		ref, err := client.Class("Files").Create(fileMap)
		if err != nil {
			t.Fatal(err)
		}

		object, err := ref.Get()
		if err != nil {
			t.Fatal(err)
		}

		if reflect.TypeOf(object.fields["file"]) != reflect.TypeOf(&File{}) {
			t.Fatal(fmt.Errorf("type of field of file shoud bd *File but %v", reflect.TypeOf(object.fields["file"])))
		}
	})
}
