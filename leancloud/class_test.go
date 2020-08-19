package leancloud

import (
	"errors"
	"testing"
)

func TestClassObject(t *testing.T) {
	client := &Client{}
	ref := client.Class("class").Object("f47ac10b58cc4372a5670e02b2c3d479")
	if ref != nil {
		if ref.c == client {
			if ref.class == "class" {
				if ref.ID == "f47ac10b58cc4372a5670e02b2c3d479" {
				} else {
					t.Fatal(errors.New("ID unmatch"))
				}
			} else {
				t.Fatal(errors.New("name of class unmatch"))
			}
		} else {
			t.Fatal(errors.New("client unmatch"))
		}
	} else {
		t.Fatal(errors.New("nil pointer of object ref"))
	}
}

func TestNewQuery(t *testing.T) {
	client := &Client{}
	query := client.Class("class").NewQuery()
	if query != nil {
		if query.c == client {
			if query.classRef.Name == "class" {
			} else {
				t.Fatal(errors.New("name of class unmatch"))
			}
		} else {
			t.Fatal(errors.New("client unmatch"))
		}
	} else {
		t.Fatal(errors.New("nil pointer of query"))
	}
}
