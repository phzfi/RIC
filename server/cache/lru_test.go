package cache

import "testing"

func TestPushPop(t *testing.T){
	lru := new(LRU)
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

