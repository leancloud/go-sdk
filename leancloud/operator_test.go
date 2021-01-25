package leancloud

import (
	"fmt"
	"reflect"
	"testing"
)

func TestOperatorEncode(t *testing.T) {
	t.Run("Increment/Decrement", func(t *testing.T) {
		ret := encode(OpIncrement(1), false)
		if !reflect.DeepEqual(ret, map[string]interface{}{
			"__op":   "Increment",
			"amount": 1,
		}) {
			t.FailNow()
		}
	})

	t.Run("Add/AddUnique/Remove", func(t *testing.T) {
		ret := encode(OpAdd([]string{"Hello", "World"}), false)
		if !reflect.DeepEqual(ret, map[string]interface{}{
			"__op":    "Add",
			"objects": []string{"Hello", "World"},
		}) {
			t.FailNow()
		}
	})

	t.Run("Delete", func(t *testing.T) {
		ret := encode(OpDelete(), false)
		if !reflect.DeepEqual(ret, map[string]interface{}{
			"__op": "Delete",
		}) {
			t.FailNow()
		}
	})
}

func TestOperatorDecode(t *testing.T) {
	t.Run("Increment/Decrement", func(t *testing.T) {
		decodedOp, err := decode(map[string]interface{}{
			"__op":   "Increment",
			"amount": 1,
		})
		if err != nil {
			t.Fatal(err)
		}
		op, ok := decodedOp.(*Op)
		if !ok {
			t.Fatal(fmt.Errorf("bad Op type: %v", reflect.TypeOf(decodedOp)))
		}

		if op.name != "Increment" && op.objects != 1 {
			t.FailNow()
		}
	})

	t.Run("Add/AddUnique/Remove", func(t *testing.T) {
		decodedOp, err := decode(map[string]interface{}{
			"__op":    "Increment",
			"objects": []string{"Hello", "World"},
		})
		if err != nil {
			t.Fatal(err)
		}
		op, ok := decodedOp.(*Op)
		if !ok {
			t.Fatal(fmt.Errorf("bad Op type: %v", reflect.TypeOf(decodedOp)))
		}

		if op.name != "Add" && reflect.DeepEqual(op.objects, []string{"Hello", "World"}) {
			t.FailNow()
		}
	})

	t.Run("Delete", func(t *testing.T) {
		decodedOp, err := decode(map[string]interface{}{
			"__op": "Delete",
		})
		if err != nil {
			t.Fatal(err)
		}
		op, ok := decodedOp.(*Op)
		if !ok {
			t.Fatal(fmt.Errorf("bad Op type: %v", reflect.TypeOf(decodedOp)))
		}

		if op.name != "Delete" {
			t.FailNow()
		}
	})
}
