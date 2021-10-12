package leancloud

import (
	"testing"
	"time"
)

func TestQueryFind(t *testing.T) {
	jake := Staff{
		Name: "Jake",
		Age:  20,
	}

	if _, err := client.Class("Staff").Create(&jake); err != nil {
		t.Fatal(err)
	}

	meeting := Meeting{
		Title:        "Team Meeting",
		Number:       1,
		Progress:     12.5,
		Host:         jake,
		Participants: []Staff{jake, jake, jake},
		Date:         time.Now(),
		Attachment:   []byte("There is nothing attachable."),
		Location:     &GeoPoint{1, 2},
	}

	if _, err := client.Class("Meeting").Create(&meeting); err != nil {
		t.Fatal(err)
	}

	ret := make([]Meeting, 10)
	if err := client.Class("Meeting").NewQuery().EqualTo("title", "Team Meeting").Include("host").Find(&ret); err != nil {
		t.Fatal(err)
	}
}

func TestQueryFirst(t *testing.T) {
	meeting := Meeting{
		Title:      "Team Meeting",
		Number:     1,
		Progress:   12.5,
		Date:       time.Now(),
		Attachment: []byte("There is nothing attachable."),
		Location:   &GeoPoint{1, 2},
	}

	if _, err := client.Class("Meeting").Create(&meeting); err != nil {
		t.Fatal(err)
	}

	ret := new(Meeting)
	if err := client.Class("Meeting").NewQuery().EqualTo("title", "Team Meeting").First(ret); err != nil {
		t.Fatal(err)
	}

	if ret.Title != meeting.Title {
		t.Fatal("dismatch title")
	}
}

func TestQueryCount(t *testing.T) {
	meeting := Meeting{
		Title:      "Team Meeting",
		Number:     1,
		Progress:   12.5,
		Date:       time.Now(),
		Attachment: []byte("There is nothing attachable."),
		Location:   &GeoPoint{1, 2},
	}

	if _, err := client.Class("Meeting").Create(&meeting); err != nil {
		t.Fatal(err)
	}

	if count, err := client.Class("Meeting").NewQuery().EqualTo("title", "Team Meeting").Count(); err != nil {
		t.Fatal(err)
	} else {
		if count < 1 {
			t.Fatal("unexpected count of results")
		}
	}
}

func TestQueryExists(t *testing.T) {
	meeting := Meeting{
		Title:      "Team Meeting",
		Number:     1,
		Progress:   12.5,
		Date:       time.Now(),
		Attachment: []byte("There is nothing attachable."),
		Location:   &GeoPoint{1, 2},
	}

	if _, err := client.Class("Meeting").Create(&meeting); err != nil {
		t.Fatal(err)
	}

	if count, err := client.Class("Meeting").NewQuery().Exists("progress").Count(); err != nil {
		t.Fatal(err)
	} else {
		if count < 1 {
			t.Fatal("unexpected count of results")
		}
	}
}

func TestQueryNotExists(t *testing.T) {
	meeting := Meeting{
		Title:      "Team Meeting",
		Number:     1,
		Date:       time.Now(),
		Attachment: []byte("There is nothing attachable."),
		Location:   &GeoPoint{1, 2},
	}

	if _, err := client.Class("Meeting").Create(&meeting); err != nil {
		t.Fatal(err)
	}

	if count, err := client.Class("Meeting").NewQuery().NotExists("progress").Count(); err != nil {
		t.Fatal(err)
	} else {
		if count < 1 {
			t.Fatal("unexpected count of results")
		}
	}
}

type room struct {
	Object
	Name    string `json:"name"`
	Meeting interface{}
}

func TestQueryMatchesQuery(t *testing.T) {
	res := []room{}
	innerQuery := client.Class("Meeting").NewQuery().EqualTo("title", "meeting1")
	client.Class("room").NewQuery().MatchesQuery("meeting", innerQuery).Find(&res)
	if len(res) < 1 || res[0].Name != "会议室1" {
		t.Fatal("unexpected count of results or wrong results")
	}
}
func TestQueryNotMatchesQuery(t *testing.T) {
	res := []room{}
	innerQuery := client.Class("Meeting").NewQuery().EqualTo("title", "meeting1")
	client.Class("room").NewQuery().NotMatchesQuery("meeting", innerQuery).Find(&res)
	if len(res) < 1 {
		t.Fatal("unexpected count of results")
	}
	for _, v := range res {
		if v.Name == "会议室1" {
			t.Fatal("wrong results")
		}
	}
}
