package leancloud

import (
	"testing"
	"time"
)

func TestQueryFind(t *testing.T) {
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

	ret := make([]Meeting, 10)
	if err := client.Class("Meeting").NewQuery().EqualTo("title", "Team Meeting").Find(&ret); err != nil {
		t.Fatal(err)
	}

	for _, v := range ret {
		if v.Title != "Team Meeting" {
			t.Fatal("dismatch title")
		}
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
		t.Log(count)
	}
}
