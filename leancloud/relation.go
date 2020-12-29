package leancloud

type Relation struct {
	key         string
	parentClass string
}

func NewRelation(key, parentClass string) *Relation {
	return &Relation{
		key:         key,
		parentClass: parentClass,
	}
}

func (relation *Relation) Add(objects ...interface{}) {

}

func (relation *Relation) Remove(objects ...interface{}) {

}
