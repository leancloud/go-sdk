package lean

import (
	"time"
)

type Object interface {
	getObjectMeta() *ObjectMeta
	ClassName() string
}

type ObjectMeta struct {
	ObjectID  string
	className string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (meta *ObjectMeta) getObjectMeta() *ObjectMeta {
	return meta
}

func ObjectID(ID string) ObjectMeta {
	return ObjectMeta{
		ObjectID: ID,
	}
}
