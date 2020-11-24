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

	ptrDate := time.Now()
	content := []byte("Meeting Content Here")
	m := Meeting{
		Title:    "Hello LeanCloud",
		TitlePtr: nil,
		TitlePtrArray: []*string{
			&(host.Name),
		},
		NullTitlePtr: nil,
		Priority:     0,
		Done:         false,
		Progress:     12.5,
		StartedAt:    &ptrDate,
		FinishedAt:   ptrDate,
		Host:         nil,
		Alternative:  *hostObject,
		Hosts: []Object{
			*hostObject,
		},
		HostsPtrArray: []*Object{
			hostObject,
		},
		HostsArrayPtr: &[]Object{
			*hostObject,
		},
		Participants: []string{"Adams", "Baker", "Clark", "Davis", "Evans", "Frank"},
		Location: GeoPoint{
			Latitude:  1.0,
			Longitude: 2.0,
		},
		Content:    content,
		ContentPtr: &content,
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

	if ret := meeting.Equal(&m); ret != 0 {
		t.Fatal(ret)
	}

	t.Log(fmt.Sprint("\n", meeting))
}
