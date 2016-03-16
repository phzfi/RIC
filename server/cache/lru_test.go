package cache

import "testing"

func TestPushPop(t *testing.T){
	lru := NewLRUPolicy()
	keys := []cacheKey{"a", "b"}
	for _, k := range keys{
		lru.Push(k)
	}
	for _, k := range keys{
		check := lru.Pop()
		if(check != k){
			t.Fatalf("LRU should have popped %s, but popped %s", k, check)
		}
	}
}

func TestVisit(t *testing.T){
	lru := NewLRUPolicy()
	keys := []cacheKey{"a", "b", "c", "d"}
	out := []cacheKey{"b", "d", "a", "c"}
	for _, k := range keys{
		lru.Push(k)
	}
	lru.Visit("a")
	lru.Visit("c")
	
	for _, k := range out{

		check := lru.Pop()
		if(check != k){
			t.Fatalf("LRU should have popped %s, but popped %s", k, check)
		}
	}
}