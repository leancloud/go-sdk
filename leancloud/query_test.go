package leancloud

import (
	"encoding/json"
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
	} else {
		t.Log(jake)
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
	} else {
		t.Log(meeting)
	}

	ret := make([]Meeting, 10)
	if err := client.Class("Meeting").NewQuery().EqualTo("title", "Team Meeting").Include("host").Find(&ret); err != nil {
		t.Fatal(err)
	} else {
		t.Log(ret)
	}

	retJSON, err := json.MarshalIndent(ret, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("\n" + string(retJSON))
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
		t.Log(count)
	}
}
