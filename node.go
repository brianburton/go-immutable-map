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
		if this.key == nil {
			return this
		} else if !equals(this.key, key) {
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

func (this *node) setValue(value Object) *node {
	if value == this.value {
		return this
	} else {
		newNode := *this
		newNode.value = value
		return &newNode
	}
}

func (this *node) setKeyAndValue(key Object, value Object) *node {
	newNode := *this
	newNode.key = key
	newNode.value = value
	return &newNode
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
