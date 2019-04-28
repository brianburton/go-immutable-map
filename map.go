package immutableMap

type Object interface{}
type HashCode uint32
type HashFunc func(Object) HashCode
type EqualsFunc func(Object, Object) bool
type Map interface {
	Assign(key Object, value Object) Map
	Get(key Object) Object
	Delete(key Object) Map
}
type MapIterator interface {
	Next() bool
	Get() (Object, Object)
}

type mapImpl struct {
	hash   HashFunc
	equals EqualsFunc
	root   *node
}

func (this *mapImpl) withRoot(newRoot *node) *mapImpl {
	newMap := *this
	newMap.root = newRoot
	return &newMap
}

func CreateMap(hash HashFunc, equals EqualsFunc) Map {
	return &mapImpl{hash: hash, equals: equals, root: emptyNode()}
}

func (this *mapImpl) Assign(key Object, value Object) Map {
	newRoot := this.root.assign(this.hash(key), key, value, this.equals)
	return this.withRoot(newRoot)
}

func (this *mapImpl) Get(key Object) Object {
	return this.root.get(this.hash(key), key, this.equals)
}

func (this *mapImpl) Delete(key Object) Map {
	newRoot := this.root.delete(this.hash(key), key, this.equals)
	if newRoot == nil {
		newRoot = emptyNode()
	}
	return this.withRoot(newRoot)
}
