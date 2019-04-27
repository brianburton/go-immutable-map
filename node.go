package immutableMap

type node struct {
	key      Object
	value    Object
	children []*node
}

func emptyNode() *node {
	return &node{}
}

func indexForHash(hashCode HashCode) int {
	return int(hashCode & 0x0f)
}

func (this *node) getChild(index int) *node {
	if this.children == nil {
		return nil
	} else {
		return this.children[index]
	}
}

func (this *node) updatedChildren(index int, child *node) []*node {
	children := make([]*node, 32)
	if this.children != nil {
		copy(children, this.children)
	}
	children[index] = child
	return children
}

func (this *node) assign(hashCode HashCode, key Object, value Object, equals EqualsFunc) *node {
	if hashCode == 0 {
		var newNode = *this
		if this.key == nil {
			newNode.key = key
			newNode.value = value
			return &newNode
		} else if equals(this.key, key) {
			newNode.value = value
			return &newNode
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

		newNode := *this
		newNode.children = this.updatedChildren(index, newChild)
		return &newNode
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

func (this *node) delete(hashCode HashCode, key Object, equals EqualsFunc) *node {
	if hashCode == 0 {
		if this.key == nil || !equals(this.key, key) {
			panic("hash collisions not yet supported")
		} else if this.childCount() == 0 {
			return nil
		} else {
			newNode := *this
			newNode.key = nil
			newNode.value = nil
			return &newNode
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
				newNode := *this
				newNode.children = this.updatedChildren(index, newChild)
				return &newNode
			}
		}
	}
}
