package immutableMap

import "fmt"

type Set interface {
	Add(key Object) Set
	Delete(key Object) Set
	Contains(key Object) bool
	Size() int
	Iterate() SetIterator
	ForEach(v SetVisitor)
	Union(s Set) Set
	Intersection(s Set) Set
	checkInvariants(report reporter)
}

type SetIterator interface {
	Next() bool
	Get() Object
}

type SetVisitor func(Object)

type setImpl struct {
	hash   HashFunc
	equals EqualsFunc
	root   *node
	size   int
}

type setIteratorImpl struct {
	state *iteratorState
	value Object
}

func keysSet(m *mapImpl) Set {
	return &setImpl{hash: m.hash, equals: m.equals, root: m.root}
}

func (this *setImpl) withRoot(newRoot *node, delta int) *setImpl {
	newSet := *this
	newSet.root = newRoot
	newSet.size += delta
	return &newSet
}

func CreateSet(hash HashFunc, equals EqualsFunc) Set {
	return &setImpl{hash: hash, equals: equals, root: emptyNode()}
}

func (this *setImpl) Add(key Object) Set {
	newRoot, delta := this.root.assign(this.hash(key), key, nil, this.equals)
	return this.withRoot(newRoot, delta)
}

func (this *setImpl) Contains(key Object) bool {
	return this.root.contains(this.hash(key), key, this.equals)
}

func (this *setImpl) Delete(key Object) Set {
	newRoot, delta := this.root.delete(this.hash(key), key, this.equals)
	if newRoot == this.root {
		return this
	} else {
		if newRoot == nil {
			newRoot = emptyNode()
		}
		return this.withRoot(newRoot, delta)
	}
}

func (this *setImpl) Size() int {
	return this.size
}

func (this *setImpl) ForEach(v SetVisitor) {
	this.root.forEach(func(value Object, _ Object) {
		v(value)
	})
}

func (this *setImpl) Union(s Set) Set {
	var larger, smaller Set
	if this.Size() > s.Size() {
		larger, smaller = this, s
	} else {
		larger, smaller = s, this
	}
	smaller.ForEach(func(v Object) {
		larger = larger.Add(v)
	})
	return larger
}

func (this *setImpl) Intersection(s Set) Set {
	var larger, smaller Set
	if this.Size() > s.Size() {
		larger, smaller = this, s
	} else {
		larger, smaller = s, this
	}
	smaller.ForEach(func(v Object) {
		if !larger.Contains(v) {
			smaller = smaller.Delete(v)
		}
	})
	return smaller
}

func (this *setImpl) checkInvariants(report reporter) {
	this.root.checkInvariants(this.hash, this.equals, 0, report)
	size := 0
	for i := this.Iterate(); i.Next(); {
		value := i.Get()
		if !this.Contains(value) {
			report(fmt.Sprintf("value from iterator not found by contains method: value=%v", value))
		}
		size++
	}
	if this.size != size {
		report(fmt.Sprintf("Size() does not match number of keys in iterator: expected=%d actual=%d", this.size, size))
	}
	i2 := this.Iterate()
	this.ForEach(func(key Object) {
		if !i2.Next() {
			report(fmt.Sprintf("Next() returned false in ForEach"))
		}
		k2 := i2.Get()
		if !this.equals(k2, key) {
			report(fmt.Sprintf("Key mismatch between ForEach and Iterate: expected=%v actual=%v", k2, key))
		}
	})
}

func (this *setImpl) Iterate() SetIterator {
	return &setIteratorImpl{state: this.root.createIteratorState(nil)}
}

func (this *setIteratorImpl) Next() bool {
	if this.state == nil {
		return false
	} else {
		this.state, this.value, _ = this.state.currentNode.next(this.state)
		return true
	}
}

func (this *setIteratorImpl) Get() Object {
	return this.value
}
