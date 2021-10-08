package leancloud

import (
	"testing"
	"time"
)

type Staff struct {
	Object
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Meeting struct {
	Object
	Title        string    `json:"title"`
	Number       int       `json:"number"`
	Progress     float64   `json:"progress"`
	Date         time.Time `json:"date"`
	Attachment   []byte    `json:"attachment"`
	Host         Staff     `json:"host"`
	Participants []Staff   `json:"participants"`
	Location     *GeoPoint `json:"location"`
}

func TestObjectCreate(t *testing.T) {
	t.Run("Struct", func(t *testing.T) {
		meeting := Meeting{
			Title:      "Team Meeting",
			Number:     1,
			Progress:   12.5,
			Date:       time.Now(),
			Attachment: []byte("There is nothing attachable."),
			Location:   &GeoPoint{1, 2},
		}

		if ref, err := testC.Class("Meeting").Create(&meeting); err != nil {
			t.Fatal(err)
		} else {
			if ref.class == "" || ref.ID == "" {
				t.FailNow()
			}
		}
	})

	t.Run("Map", func(t *testing.T) {
		meeting := map[string]interface{}{
			"title":      "Team Meeting",
			"number":     1,
			"progress":   12.5,
			"date":       time.Now(),
			"attachment": []byte("There is nothing attachable."),
			"location":   &GeoPoint{1, 2},
		}

		if ref, err := testC.Class("Meeting").Create(meeting); err != nil {
			t.Fatal(err)
		} else {
			if ref.class == "" || ref.ID == "" {
				t.FailNow()
			}
		}
	})
}

func TestObjectGet(t *testing.T) {
	t.Run("Custom", func(t *testing.T) {
		jake := Staff{
			Name: "Jake",
			Age:  20,
		}

		_, err := testC.Class("Staff").Create(&jake)
		if err != nil {
			t.Fatal()
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

		_, err = testC.Class("Meeting").Create(&meeting)
		if err != nil {
			t.Fatal(err)
		}

		newMeeting := new(Meeting)
		if err := testC.Class("Meeting").ID(meeting.ID).Get(newMeeting); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Bare", func(t *testing.T) {
		meeting := map[string]interface{}{
			"title":      "Team Meeting",
			"number":     1,
			"progress":   12.5,
			"date":       time.Now(),
			"attachment": []byte("There is nothing attachable."),
			"location":   &GeoPoint{1, 2},
		}

		ref, err := testC.Class("Meeting").Create(meeting)
		if err != nil {
			t.Fatal(err)
		}

		object := new(Object)
		if err := testC.Class("Meeting").ID(ref.ID).Get(object); err != nil {
			t.Fatal(err)
		}
	})
}

func TestObjectSet(t *testing.T) {
	meeting := Meeting{
		Title:      "Team Meeting",
		Number:     1,
		Progress:   12.5,
		Date:       time.Now(),
		Attachment: []byte("There is nothing attachable."),
		Location:   &GeoPoint{1, 2},
	}

	if _, err := testC.Class("Meeting").Create(&meeting); err != nil {
		t.Fatal(err)
	}

	if err := testC.Object(meeting).Set("title", "Another Team Meeting"); err != nil {
		t.Fatal(err)
	}

	newMeeting := new(Meeting)
	if err := testC.Object(meeting).Get(newMeeting); err != nil {
		t.Fatal(err)
	}
}

func TestObjectUpdate(t *testing.T) {
	t.Run("Struct", func(t *testing.T) {
		meeting := Meeting{
			Title:      "Team Meeting",
			Number:     1,
			Progress:   12.5,
			Date:       time.Now(),
			Attachment: []byte("There is nothing attachable."),
			Location:   &GeoPoint{1, 2},
		}

		if _, err := testC.Class("Meeting").Create(&meeting); err != nil {
			t.Fatal(err)
		}

		diff := &Meeting{
			Title:    "Another Team Meeting",
			Number:   2,
			Progress: 13.5,
			Date:     time.Now(),
		}

		if err := testC.Object(meeting).Update(diff); err != nil {
			t.Fatal(err)
		}

		newMeeting := new(Meeting)
		if err := testC.Object(meeting).Get(newMeeting); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Map", func(t *testing.T) {
		meeting := map[string]interface{}{
			"title":      "Team Meeting",
			"number":     1,
			"progress":   12.5,
			"date":       time.Now(),
			"attachment": []byte("There is nothing attachable."),
			"location":   &GeoPoint{1, 2},
		}

		ref, err := testC.Class("Meeting").Create(meeting)
		if err != nil {
			t.Fatal(err)
		}

		if err := testC.Class("Meeting").ID(ref.ID).Update(map[string]interface{}{
			"title":    "Another Team Meeting",
			"number":   2,
			"progress": 13.5,
			"date":     time.Now(),
		}); err != nil {
			t.Fatal(err)
		}

		newMeeting := new(Meeting)
		if err := testC.Class("Meeting").ID(ref.ID).Get(newMeeting); err != nil {
			t.Fatal(err)
		}
	})
}

func TestObjectUpdateWithQuery(t *testing.T) {
	meeting := Meeting{
		Title:      "Team Meeting",
		Number:     1,
		Progress:   13.5,
		Date:       time.Now(),
		Attachment: []byte("There is nothing attachable."),
		Location:   &GeoPoint{1, 2},
	}

	if _, err := testC.Class("Meeting").Create(&meeting); err != nil {
		t.Fatal(err)
	}

	diff := &Meeting{
		Title:    "Another Team Meeting",
		Number:   2,
		Progress: 14.5,
		Date:     time.Now(),
	}

	if err := testC.Object(meeting).UpdateWithQuery(diff, testC.Class("Meeting").NewQuery().EqualTo("progress", 13.5)); err != nil {
		t.Fatal(err)
	}

	newMeeting := new(Meeting)
	if err := testC.Object(meeting).Get(newMeeting); err != nil {
		t.Fatal(err)
	}
}

func TestObjectDestroy(t *testing.T) {
	meeting := Meeting{
		Title:      "Team Meeting",
		Number:     1,
		Progress:   12.5,
		Date:       time.Now(),
		Attachment: []byte("There is nothing attachable."),
		Location:   &GeoPoint{1, 2},
	}

	if _, err := testC.Class("Meeting").Create(&meeting); err != nil {
		t.Fatal(err)
	}

	if err := testC.Object(meeting).Destroy(UseMasterKey(true)); err != nil {
		t.Fatal(err)
	}

	newMeeting := new(Meeting)
	if err := testC.Object(meeting).Get(newMeeting); err == nil {
		if newMeeting.ID != "" {
			t.FailNow()
		}
	}
}
