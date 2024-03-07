package cache

import "testing"

func Test_CacheHappyCase(t *testing.T) {
	cache := New[int]()
	inValue := 1
	cache.Set("test", inValue)
	value, ok := cache.Get("test")
	if !ok {
		t.Error("set value not in cache")
	}
	if value != inValue {
		t.Errorf("expected %d got %d", inValue, value)
	}
}

func Test_CacheEmptyCase(t *testing.T) {
	cache := New[int]()
	_, ok := cache.Get("test")
	if ok {
		t.Error("expected false for an unset key")
	}
}
