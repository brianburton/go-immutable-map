package immutableMap

type node struct {
	key      Object
	value    Object
	children []*node
}

func emptyNode() *node {
	return &node{}
}

func (this *node) assign(hashCode HashCode, key Object, value Object, equals EqualsFunc) *node {
	if hashCode == 0 {
		if this.key == nil {
			return this.setKeyAndValue(key, value)
		} else if equals(this.key, key) {
			return this.setValue(value)
		} else {
			panic("hash collisions not yet supported")
		}
	} else {
		index := indexForHash(hashCode)
		oldChild := this.getChild(index)
		if oldChild == nil {
			oldChild = emptyNode()
		}
		newChild := oldChild.assign(hashCode>>5, key, value, equals)
		return this.setChild(index, newChild)
	}
}

func (this *node) get(hashCode HashCode, key Object, equals EqualsFunc) Object {
	if hashCode == 0 {
		if this.key == nil || !equals(this.key, key) {
			panic("hash collisions not yet supported")
		} else {
			return this.value
		}
	} else {
		index := indexForHash(hashCode)
		oldChild := this.getChild(index)
		if oldChild == nil {
			return nil
		} else {
			return oldChild.get(hashCode>>5, key, equals)
		}
	}
}

func (this *node) delete(hashCode HashCode, key Object, equals EqualsFunc) *node {
	if hashCode == 0 {
		if this.key == nil || !equals(this.key, key) {
			panic("hash collisions not yet supported")
		} else if this.childCount() == 0 {
			return nil
		} else {
			return this.setKeyAndValue(nil, nil)
		}
	} else {
		index := indexForHash(hashCode)
		oldChild := this.getChild(index)
		if oldChild == nil {
			return this
		} else {
			newChild := oldChild.delete(hashCode>>5, key, equals)
			if newChild == nil && this.key == nil && this.childCount() == 1 {
				return nil
			} else {
				return this.setChild(index, newChild)
			}
		}
	}
}

func indexForHash(hashCode HashCode) int {
	return int(hashCode & 0x0f)
}

func (this *node) setValue(value Object) *node {
	newNode := *this
	newNode.value = value
	return &newNode
}

func (this *node) setKeyAndValue(key Object, value Object) *node {
	newNode := *this
	newNode.key = key
	newNode.value = value
	return &newNode
}

func (this *node) childCount() int {
	if this.children == nil {
		return 0
	}
	count := 0
	for _, c := range this.children {
		if c != nil {
			count++
		}
	}
	return count
}

func (this *node) getChild(index int) *node {
	if this.children == nil {
		return nil
	} else {
		return this.children[index]
	}
}

func (this *node) setChild(index int, child *node) *node {
	newNode := *this
	newNode.children = make([]*node, 32)
	if this.children != nil {
		copy(newNode.children, this.children)
	}
	newNode.children[index] = child
	return &newNode
}
