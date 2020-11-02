package leancloud

import (
	"fmt"
	"testing"
	"time"
)

func TestQueryFind(t *testing.T) {
	data, err := beforeTestQuery(c.Class("Todo").NewQuery())
	if err != nil {
		t.Fatal(err)
	}

	todo := data.(*Todo)

	t.Run("EqualTo", func(t *testing.T) {
		results, err := c.Class("Todo").NewQuery().EqualTo("title", "Team Meeting").Find()
		if err != nil {
			t.Fatal(err)
		}
		for _, v := range results {
			if v.fields["title"].(string) != todo.Title {
				t.Fatal(fmt.Errorf("wrong result of querying"))
			}
		}
	})

	t.Run("NotEqualTo", func(t *testing.T) {
		results, err := c.Class("Todo").NewQuery().NotEqualTo("done", false).Find()
		if err != nil {
			t.Fatal(err)
		}

		if len(results) == 0 {
			t.Fatal(fmt.Errorf("wrong result of querying"))
		}
	})

	t.Run("GreaterThan", func(t *testing.T) {
		results, err := c.Class("Todo").NewQuery().GreaterThan("priority", 10).Find()
		if err != nil {
			t.Fatal(err)
		}

		if len(results) != 0 {
			t.Fatal(fmt.Errorf("wrong result of querying"))
		}
	})

	t.Run("GreaterThanOrEqualTo", func(t *testing.T) {
		results, err := c.Class("Todo").NewQuery().GreaterThanOrEqualTo("priority", 11).Find()
		if err != nil {
			t.Fatal(err)
		}

		if len(results) != 0 {
			t.Fatal(fmt.Errorf("wrong result of querying"))
		}
	})

	t.Run("LessThan", func(t *testing.T) {
		results, err := c.Class("Todo").NewQuery().LessThan("priority", 1).Find()
		if err != nil {
			t.Fatal(err)
		}

		if len(results) != 0 {
			t.Fatal(fmt.Errorf("wrong result of querying"))
		}
	})

	t.Run("LessThanOrEqualTo", func(t *testing.T) {
		results, err := c.Class("Todo").NewQuery().LessThanOrEqualTo("priority", 10).Find()
		if err != nil {
			t.Fatal(err)
		}

		if len(results) == 0 {
			t.Fatal(fmt.Errorf("wrong result of querying"))
		}
	})
}

func TestQueryFirst(t *testing.T) {
	result, err := c.Class("Todo").NewQuery().LessThan("priority", 10).First()
	if err != nil {
		t.Fatal(err)
	}

	if int(result.fields["priority"].(float64)) >= 10 {
		t.Fatal(fmt.Errorf("wrong result of querying"))
	}
}

func TestQueryCount(t *testing.T) {
	count, err := c.Class("Todo").NewQuery().GreaterThan("priority", 1).Count()
	if err != nil {
		t.Fatal(err)
	}

	if count == 0 {
		t.Fatal(fmt.Errorf("wrong result of querying"))
	}
}

func beforeTestQuery(query *Query) (interface{}, error) {
	todo := Todo{
		Title:      "Team Meeting",
		Priority:   1,
		Done:       false,
		Progress:   12.5,
		FinishedAt: time.Now(),
	}

	if _, err := c.Class("Todo").Create(todo); err != nil {
		return nil, err
	}

	return &todo, nil
}
