package immutableMap

type Set interface {
	Add(key Object) Set
	Delete(key Object) Set
	Contains(key Object) bool
	Iterate() SetIterator
}

type SetIterator interface {
	Next() bool
	Get() Object
}

type setImpl struct {
	hash   HashFunc
	equals EqualsFunc
	root   *node
}

type setIteratorImpl struct {
	state *iteratorState
	value Object
}

func (this *setImpl) withRoot(newRoot *node) *setImpl {
	newSet := *this
	newSet.root = newRoot
	return &newSet
}

func CreateSet(hash HashFunc, equals EqualsFunc) Set {
	return &setImpl{hash: hash, equals: equals, root: emptyNode()}
}

func (this *setImpl) Add(key Object) Set {
	newRoot := this.root.assign(this.hash(key), key, nil, this.equals)
	return this.withRoot(newRoot)
}

func (this *setImpl) Contains(key Object) bool {
	return this.root.contains(this.hash(key), key, this.equals)
}

func (this *setImpl) Delete(key Object) Set {
	newRoot := this.root.delete(this.hash(key), key, this.equals)
	if newRoot == this.root {
		return this
	} else {
		if newRoot == nil {
			newRoot = emptyNode()
		}
		return this.withRoot(newRoot)
	}
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
