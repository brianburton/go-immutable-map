package immutableMap

import (
	"fmt"
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

func val(index int) string {
	return fmt.Sprintf("%v", index)
}

func TestMap(t *testing.T) {
	m := Create(stringHash, stringEquals)
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
