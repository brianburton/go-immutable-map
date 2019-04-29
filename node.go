package immutableMap

import (
	"math/bits"
)

type node struct {
	key      Object
	value    Object
	bitmask  uint32
	children []*node
}

type iteratorState struct {
	next         *iteratorState
	currentNode  *node
	currentIndex int
}

func (this *node) isEmpty() bool {
	return this.key == nil && this.bitmask == 0
}

func emptyNode() *node {
	return &node{}
}

func (this *node) assign(hashCode HashCode, key Object, value Object, equals EqualsFunc) *node {
	if hashCode == 0 {
		if this.key == nil {
			return this.setKeyAndValue(key, value)
		} else if equals(this.key, key) {
			return this.setKeyAndValue(key, value)
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
		if newChild == oldChild {
			return this
		} else {
			return this.setChild(index, newChild)
		}
	}
}

func (this *node) get(hashCode HashCode, key Object, equals EqualsFunc) Object {
	if hashCode == 0 {
		if this.key == nil || !equals(this.key, key) {
			return nil
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

func (this *node) contains(hashCode HashCode, key Object, equals EqualsFunc) bool {
	if hashCode == 0 {
		return this.key != nil && equals(this.key, key)
	} else {
		index := indexForHash(hashCode)
		child := this.getChild(index)
		return child != nil && child.contains(hashCode>>5, key, equals)
	}
}

func (this *node) delete(hashCode HashCode, key Object, equals EqualsFunc) *node {
	if hashCode == 0 {
		return this.deleteKey(key, equals)
	} else {
		index := indexForHash(hashCode)
		oldChild := this.getChild(index)
		if oldChild == nil {
			return this
		} else {
			newChild := oldChild.delete(hashCode>>5, key, equals)
			if newChild == oldChild {
				return this
			} else if newChild == nil {
				return this.deleteChild(index)
			} else {
				return this.setChild(index, newChild)
			}
		}
	}
}

func indexForHash(hashCode HashCode) int {
	return int(hashCode & 0x0f)
}

func (this *node) setKeyAndValue(key Object, value Object) *node {
	if this.key == key && this.value == value {
		return this
	} else {
		newNode := *this
		newNode.key = key
		newNode.value = value
		return &newNode
	}
}

func (this *node) deleteKey(key Object, equals EqualsFunc) *node {
	if this.key == nil {
		return this
	} else if !equals(this.key, key) {
		panic("hash collisions not yet supported")
	} else if this.childCount() == 0 {
		return nil
	} else {
		newNode := *this
		newNode.key = nil
		newNode.value = nil
		return &newNode
	}
}

func (this *node) childCount() int {
	return bits.OnesCount32(this.bitmask)
}

func (this *node) getChild(index int) *node {
	indexBit := indexBit(index)
	if this.bitmask&indexBit == 0 {
		return nil
	} else {
		realIndex := this.realIndex(indexBit)
		return this.children[realIndex]
	}
}

func (this *node) realIndex(indexBit uint32) int {
	trailingBits := indexBit - 1
	realIndex := bits.OnesCount32(this.bitmask & trailingBits)
	return realIndex
}

func indexBit(index int) uint32 {
	var indexBit uint32 = 1 << uint32(index)
	return indexBit
}

func (this *node) setChild(index int, child *node) *node {
	newNode := *this
	indexBit := indexBit(index)
	if this.children == nil {
		newNode.children = make([]*node, 1)
		newNode.children[0] = child
		newNode.bitmask = indexBit
	} else {
		realIndex := this.realIndex(indexBit)
		if this.bitmask&indexBit != 0 {
			newNode.children = make([]*node, len(this.children))
			copy(newNode.children, this.children)
			newNode.children[realIndex] = child
		} else {
			newNode.children = make([]*node, len(this.children)+1)
			copy(newNode.children, this.children[0:realIndex])
			newNode.children[realIndex] = child
			copy(newNode.children[realIndex+1:], this.children[realIndex:])
			newNode.bitmask |= indexBit
		}
	}
	return &newNode
}

func (this *node) deleteChild(index int) *node {
	indexBit := indexBit(index)
	if this.bitmask&indexBit == 0 {
		panic("attempting to delete non-existent child")
	}

	newNode := *this
	if this.childCount() == 1 {
		if this.key == nil {
			return nil
		} else {
			newNode.children = nil
			newNode.bitmask = 0
		}
	} else {
		realIndex := this.realIndex(indexBit)
		newNode.children = make([]*node, len(this.children)-1)
		copy(newNode.children, this.children[0:realIndex])
		copy(newNode.children[realIndex:], this.children[realIndex+1:])
		newNode.bitmask &= ^indexBit
	}
	return &newNode
}

func (this *node) createIteratorState(nextState *iteratorState) *iteratorState {
	if this.isEmpty() {
		return nextState
	} else {
		var startingIndex int
		if this.key == nil {
			startingIndex = 0
		} else {
			startingIndex = -1
		}
		return &iteratorState{currentNode: this, next: nextState, currentIndex: startingIndex}
	}
}

func (this *node) next(state *iteratorState) (*iteratorState, Object, Object) {
	if state == nil || state.currentNode != this {
		state = this.createIteratorState(state)
	}
	if state.currentIndex == -1 {
		state.currentIndex++
		if len(this.children) > 0 {
			return state, this.key, this.value
		} else {
			return state.next, this.key, this.value
		}
	}
	child := this.children[state.currentIndex]
	state.currentIndex++
	if state.currentIndex == len(this.children) {
		return child.next(state.next)
	} else {
		return child.next(state)
	}
}
