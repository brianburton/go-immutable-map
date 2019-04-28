package immutableMap

import (
	"fmt"
	"strconv"
	"testing"
)

func stringEquals(a Object, b Object) bool {
	return a.(string) == b.(string)
}

func stringHash(a Object) HashCode {
	var val HashCode = 0
	for _, c := range a.(string) {
		val = 31*val + HashCode(c)
	}
	return val
}

func digitsHash(a Object) HashCode {
	i, _ := strconv.Atoi(a.(string))
	return HashCode(i)
}

func val(index int) string {
	return fmt.Sprintf("%v", index)
}

func TestMap(t *testing.T) {
	m := CreateMap(stringHash, stringEquals)
	for i := -2000; i <= 2000; i++ {
		key := val(i)
		m = m.Assign(key, i)
	}

	m = m.Assign(val(0), -1)
	m = m.Assign(val(0), 0)

	for i := 2000; i >= -2000; i-- {
		key := val(i)
		v := m.Get(key)
		if v.(int) != i {
			t.Error(fmt.Sprintf("expected %v but got %v for key %v", i, v, key))
		}
	}

	for i := -2000; i <= 0; i++ {
		key := val(i)
		m = m.Delete(key)
	}

	for i := 2000; i >= -5; i-- {
		key := val(i)
		m = m.Delete(key)
	}

	for i := 2000; i >= -2000; i-- {
		key := val(i)
		v := m.Get(key)
		if v != nil {
			t.Error(fmt.Sprintf("expected nil but got %v for key %v", v, key))
		}
	}

	if m == nil {
		t.Error(fmt.Sprintf("can't really happen"))
	}
}

func TestSet(t *testing.T) {
	s := CreateSet(stringHash, stringEquals)
	for i := -2000; i <= 2000; i++ {
		key := val(i)
		s = s.Add(key)
	}

	s = s.Add(val(0))

	for i := 2000; i >= -2000; i-- {
		key := val(i)
		v := s.Contains(key)
		if !v {
			t.Error(fmt.Sprintf("expected true but got %v for key %v", v, key))
		}
	}

	for i := -2000; i <= 0; i++ {
		key := val(i)
		s = s.Delete(key)
	}

	for i := 2000; i >= -5; i-- {
		key := val(i)
		s = s.Delete(key)
	}

	for i := 2000; i >= -2000; i-- {
		key := val(i)
		v := s.Contains(key)
		if v {
			t.Error(fmt.Sprintf("expected false but got %v for key %v", v, key))
		}
	}

	if s == nil {
		t.Error(fmt.Sprintf("can't really happen"))
	}
}

func TestMapPaths(t *testing.T) {
	m := CreateMap(digitsHash, stringEquals)

	m = m.Assign(keyForPath([]int{1, 2, 3}), 3)

	m = m.Delete(keyForPath([]int{1, 2, 3}))

	m = m.Assign(keyForPath([]int{1, 1, 1}), 1)
	m = m.Assign(keyForPath([]int{1, 2, 2}), 2)
	m = m.Assign(keyForPath([]int{2, 3, 3}), 3)

	if v := m.Get(keyForPath([]int{1, 1})); v != nil {
		t.Error(fmt.Sprintf("Get returned %v for key %v", v, keyForPath([]int{1, 1})))
	}

	m.Delete(keyForPath([]int{1, 1}))

	m = m.Delete(keyForPath([]int{1, 2, 2}))
	m = m.Delete(keyForPath([]int{1, 1, 1}))
	m = m.Delete(keyForPath([]int{2, 3, 3}))

	m = nil
}

func TestSetPaths(t *testing.T) {
	s := CreateSet(digitsHash, stringEquals)

	s = s.Add(keyForPath([]int{1, 2, 3}))

	s = s.Delete(keyForPath([]int{1, 2, 3}))

	s = s.Add(keyForPath([]int{1, 1, 1}))
	s = s.Add(keyForPath([]int{1, 2, 2}))
	s = s.Add(keyForPath([]int{2, 3, 3}))

	if s.Contains(keyForPath([]int{1, 1})) {
		t.Error(fmt.Sprintf("Contains returned true for key %s", keyForPath([]int{1, 1})))
	}

	del := s.Delete(keyForPath([]int{1, 1}))
	if del != s {
		t.Error(fmt.Sprintf("Delete returned new set for key %s", keyForPath([]int{1, 1})))
	}

	s = s.Delete(keyForPath([]int{1, 2, 2}))
	s = s.Delete(keyForPath([]int{1, 1, 1}))
	s = s.Delete(keyForPath([]int{2, 3, 3}))

	s = nil
}

func keyForPath(indexes []int) string {
	key := 0
	for m, i := range indexes {
		key += i << (5 * uint(m))
	}
	return val(key)
}
