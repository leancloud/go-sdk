package leancloud

import (
	"fmt"
	"testing"
	"time"
)

func TestObjectToStruct(t *testing.T) {
	host := Host{
		Name:       "Chris",
		Department: "Engineering",
		Leader:     "Wu",
		Level:      3,
	}
	hostRef, err := c.Class("Host").Create(host)
	if err != nil {
		t.Fatal(err)
	}

	hostObject, err := hostRef.Get()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(hostObject)
	ptrDate := time.Now()
	m := Meeting{
		Title:        "",
		Priority:     0,
		Done:         false,
		Progress:     12.5,
		StartedAt:    &ptrDate,
		FinishedAt:   ptrDate,
		Host:         hostObject,
		Participants: []string{"Adams", "Baker", "Clark", "Davis", "Evans", "Frank"},
		Location: GeoPoint{
			Latitude:  1.0,
			Longitude: 2.0,
		},
		//Content: []byte("Meeting Content Here"),
	}

	ref, err := c.Class("Meetings").Create(m)
	if err != nil {
		t.Fatal(err)
	}

	mObject, err := ref.Get()
	if err != nil {
		t.Fatal(err)
	}

	meeting := new(Meeting)
	if err := mObject.ToStruct(meeting); err != nil {
		t.Fatal(err)
	}

	if !meeting.Equal(&m) {
		t.FailNow()
	}

	t.Log(fmt.Sprint("\n", meeting))
}
