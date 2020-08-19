package leancloud

import (
	"errors"
	"testing"
)

func TestClassObject(t *testing.T) {
	client := &Client{}
	ref := client.Class("class").Object("f47ac10b58cc4372a5670e02b2c3d479")

	if ref == nil {
		t.Fatal(errors.New("nil pointer of object ref"))
	}
	if ref.c != client {
		t.Fatal(errors.New("client unmatch"))
	}
	if ref.class != "class" {
		t.Fatal(errors.New("name of class unmatch"))
	}
	if ref.ID != "f47ac10b58cc4372a5670e02b2c3d479" {
		t.Fatal(errors.New("ID unmatch"))
	}
}

func TestNewQuery(t *testing.T) {
	client := &Client{}
	query := client.Class("class").NewQuery()

	if query == nil {
		t.Fatal(errors.New("nil pointer of query"))
	}
	if query.c != client {
		t.Fatal(errors.New("client unmatch"))
	}
	if query.classRef.Name != "class" {
		t.Fatal(errors.New("name of class unmatch"))
	}
}
