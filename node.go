package immutableMap

import (
	"fmt"
	"math/bits"
)

type keyValueList struct {
	next  *keyValueList
	key   Object
	value Object
}

type node struct {
	keys     *keyValueList
	bitmask  uint32
	children []*node
}

type iteratorState struct {
	next         *iteratorState
	currentNode  *node
	currentIndex int
	currentKey   *keyValueList
}

func (this *node) isEmpty() bool {
	return this.keys == nil && this.bitmask == 0
}

func emptyNode() *node {
	return &node{}
}

func (this *node) assign(hashCode HashCode, key Object, value Object, equals EqualsFunc) (*node, int) {
	if hashCode == 0 {
		return this.setKeyAndValue(key, value, equals)
	} else {
		index := indexForHash(hashCode)
		oldChild := this.getChild(index)
		if oldChild == nil {
			oldChild = emptyNode()
		}
		newChild, delta := oldChild.assign(hashCode>>5, key, value, equals)
		if newChild == oldChild {
			return this, delta
		} else {
			return this.setChild(index, newChild), delta
		}
	}
}

func (this *node) get(hashCode HashCode, key Object, equals EqualsFunc) Object {
	if hashCode == 0 {
		return this.getValueForKey(key, equals)
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
		return this.containsValueForKey(key, equals)
	} else {
		index := indexForHash(hashCode)
		child := this.getChild(index)
		return child != nil && child.contains(hashCode>>5, key, equals)
	}
}

func (this *node) delete(hashCode HashCode, key Object, equals EqualsFunc) (*node, int) {
	if hashCode == 0 {
		return this.deleteKey(key, equals)
	} else {
		index := indexForHash(hashCode)
		oldChild := this.getChild(index)
		if oldChild == nil {
			return this, 0
		} else {
			newChild, delta := oldChild.delete(hashCode>>5, key, equals)
			if newChild == oldChild {
				return this, 0
			} else if newChild == nil {
				return this.deleteChild(index), delta
			} else {
				return this.setChild(index, newChild), delta
			}
		}
	}
}

func indexForHash(hashCode HashCode) int {
	return int(hashCode & 0x0f)
}

func (this *node) containsValueForKey(key Object, equals EqualsFunc) bool {
	for kvp := this.keys; kvp != nil; kvp = kvp.next {
		if equals(key, kvp.key) {
			return true
		}
	}
	return false
}

func (this *node) getValueForKey(key Object, equals EqualsFunc) Object {
	for kvp := this.keys; kvp != nil; kvp = kvp.next {
		if equals(key, kvp.key) {
			return kvp.value
		}
	}
	return nil
}

func (this *node) setKeyAndValue(key Object, value Object, equals EqualsFunc) (*node, int) {
	var newKeys *keyValueList
	delta := 0
	if this.keys == nil {
		newKeys = &keyValueList{key: key, value: value}
		delta = 1
	} else {
		changed := false
		for kvp := this.keys; kvp != nil; kvp = kvp.next {
			if equals(kvp.key, key) {
				if kvp.value == value {
					return this, 0
				}
				newKeys = &keyValueList{key: key, value: value, next: newKeys}
				changed = true
			} else {
				newKeys = &keyValueList{key: kvp.key, value: kvp.value, next: newKeys}
			}
		}
		if !changed {
			newKeys = &keyValueList{key: key, value: value, next: this.keys}
			delta = 1
		}
	}
	newNode := *this
	newNode.keys = newKeys
	return &newNode, delta
}

func (this *node) deleteKey(key Object, equals EqualsFunc) (*node, int) {
	if this.keys == nil {
		return this, 0
	}

	changed := false
	var newKeys *keyValueList
	for kvp := this.keys; kvp != nil; kvp = kvp.next {
		if equals(kvp.key, key) {
			changed = true
		} else {
			newKeys = &keyValueList{key: kvp.key, value: kvp.value, next: newKeys}
		}
	}
	if !changed {
		return this, 0
	} else if newKeys == nil && this.childCount() == 0 {
		return nil, -1
	} else {
		newNode := *this
		newNode.keys = newKeys
		return &newNode, -1
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
	newNode := *this
	if this.childCount() == 1 {
		if this.keys == nil {
			return nil
		} else {
			newNode.children = nil
			newNode.bitmask = 0
		}
	} else {
		indexBit := indexBit(index)
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
		if this.keys == nil {
			startingIndex = 0
		} else {
			startingIndex = -1
		}
		return &iteratorState{next: nextState, currentNode: this, currentIndex: startingIndex, currentKey: this.keys}
	}
}

func (this *node) next(state *iteratorState) (*iteratorState, Object, Object) {
	if state == nil || state.currentNode != this {
		state = this.createIteratorState(state)
	}
	if state.currentIndex == -1 {
		kvp := state.currentKey
		state.currentKey = kvp.next
		if state.currentKey != nil {
			return state, kvp.key, kvp.value
		} else {
			state.currentIndex++
			if len(this.children) > 0 {
				return state, kvp.key, kvp.value
			} else {
				return state.next, kvp.key, kvp.value
			}
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

func (this *node) checkInvariants(hash HashFunc, equals EqualsFunc, shift uint, report reporter) {
	for kvp := this.keys; kvp != nil; kvp = kvp.next {
		for other := kvp.next; other != nil; other = other.next {
			if equals(kvp.key, other.key) {
				report(fmt.Sprintf("duplicate key detected: key=%v", kvp.key))
			}
		}
		if shiftedHash := hash(kvp.key) >> shift; shiftedHash != 0 {
			report(fmt.Sprintf("key with non-zero hash detected: key=%v shiftedHash=%d", kvp.key, shiftedHash))
		}
	}

	if this.bitmask != 0 && this.children == nil {
		report("nil children with non-zero bitmask")
	} else if this.bitmask == 0 && this.children != nil {
		report("non-nil children with zero bitmask")
	}

	if this.children != nil {
		if bitsLength := bits.OnesCount32(this.bitmask); bitsLength != len(this.children) {
			report(fmt.Sprintf("bitmask count differs from children length: bitmask=%x bitsLength=%d sliceLength=%d", this.bitmask, bitsLength, len(this.children)))
		}
		for _, c := range this.children {
			c.checkInvariants(hash, equals, shift+5, report)
		}
	}
}
