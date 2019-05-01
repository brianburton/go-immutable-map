package immutableMap

import "fmt"

type Object interface{}
type HashCode uint32
type HashFunc func(Object) HashCode
type EqualsFunc func(Object, Object) bool
type reporter func(message string)

type Map interface {
	Assign(key Object, value Object) Map
	Get(key Object) Object
	Delete(key Object) Map
	Keys() Set
	Size() int
	Iterate() MapIterator
	checkInvariants(report reporter)
}
type MapIterator interface {
	Next() bool
	Get() (Object, Object)
}

type mapImpl struct {
	hash   HashFunc
	equals EqualsFunc
	root   *node
	size   int
}

type mapIteratorImpl struct {
	state *iteratorState
	key   Object
	value Object
}

func (this *mapImpl) withRoot(newRoot *node, delta int) *mapImpl {
	newMap := *this
	newMap.root = newRoot
	newMap.size += delta
	return &newMap
}

func CreateMap(hash HashFunc, equals EqualsFunc) Map {
	return &mapImpl{hash: hash, equals: equals, root: emptyNode()}
}

func (this *mapImpl) Assign(key Object, value Object) Map {
	newRoot, delta := this.root.assign(this.hash(key), key, value, this.equals)
	return this.withRoot(newRoot, delta)
}

func (this *mapImpl) Get(key Object) Object {
	return this.root.get(this.hash(key), key, this.equals)
}

func (this *mapImpl) Delete(key Object) Map {
	newRoot, delta := this.root.delete(this.hash(key), key, this.equals)
	if newRoot == nil {
		newRoot = emptyNode()
	}
	return this.withRoot(newRoot, delta)
}

func (this *mapImpl) Keys() Set {
	return keysSet(this)
}

func (this *mapImpl) Size() int {
	return this.size
}

func (this *mapImpl) Iterate() MapIterator {
	return &mapIteratorImpl{state: this.root.createIteratorState(nil)}
}

func (this *mapImpl) checkInvariants(report reporter) {
	this.root.checkInvariants(this.hash, this.equals, 0, report)
	size := 0
	for i := this.Iterate(); i.Next(); {
		key, expected := i.Get()
		actual := this.Get(key)
		if expected != actual {
			report(fmt.Sprintf("Get returned incorrect result: key=%v expected=%v actual=%v", key, expected, actual))
		}
		size++
	}
	if this.size != size {
		report(fmt.Sprintf("Size() does not match number of keys in iterator: expected=%d actual=%d", this.size, size))
	}
}

func (this *mapIteratorImpl) Next() bool {
	if this.state == nil {
		return false
	} else {
		this.state, this.key, this.value = this.state.currentNode.next(this.state)
		return true
	}
}

func (this *mapIteratorImpl) Get() (Object, Object) {
	return this.key, this.value
}
