package cache

import "testing"

func TestPushPop(t *testing.T) {
	lru := NewLRU()
	keys := []string{"a", "b"}
	for _, k := range keys {
		lru.Push(k)
	}
	for _, k := range keys {
		check := lru.Pop()
		if check != k {
			t.Fatalf("LRU should have popped %s, but popped %s", k, check)
		}
	}
}

func TestVisit(t *testing.T) {
	lru := NewLRU()
	keys := []string{"a", "b", "c", "d"}
	out := []string{"b", "d", "a", "c"}
	for _, k := range keys {
		lru.Push(k)
	}
	lru.Visit("a")
	lru.Visit("c")

	for _, k := range out {

		check := lru.Pop()
		if check != k {
			t.Fatalf("LRU should have popped %s, but popped %s", k, check)
		}
	}
}

func TestUnderflow(t *testing.T) {
	lru := NewLRU()

	defer func() {
		if recover() == nil {
			t.Fatal("No panic occurred.")
		}
	}()

	lru.Pop()
}

func TestKeyCollision(t *testing.T) {
	lru := NewLRU()

	lru.Push("asd")
	lru.Push("asd")

	lru.Pop()

	lru.Visit("asd")
}
