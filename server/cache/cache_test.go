package cache

import (
	"testing"
	"github.com/phzfi/RIC/server/ops"
)

const (
	Visit = iota
	Push
	Pop
)

type Log map[cacheKey][]uint

type DummyPolicy struct{
	fifo Policy
	
	loki Log
	pops int
}

func (d DummyPolicy) Visit(k cacheKey) {
	d.log(k, Visit)
	d.fifo.Visit(k)
}

func (d DummyPolicy) log(k cacheKey, t uint){
	d.loki[k] = append(d.loki[k], t)
}

func (d DummyPolicy) Push(k cacheKey) {
	d.log(k, Push)
	d.fifo.Push(k)
}

func (d *DummyPolicy) Pop() cacheKey {
	d.pops += 1
	return d.fifo.Pop()
}

func NewDummyPolicy(log Log) *DummyPolicy{
	return &DummyPolicy{fifo: &FIFO{}, loki: log}
}

func setup() (dp *DummyPolicy, cache *Cache){
	dp = NewDummyPolicy(make(Log))
	cache = NewCache(dp, 100)
	return
}

func TestCache(t *testing.T) {
	id := []ops.Operation{&DummyOperation{}}
	dp, cache := setup()
	
	found := func()bool{
		_, ok := cache.GetBlob(id)
		return ok
	}

	if found(){
		t.Fatal("Cache claimed to contain a blob that was never added")
	}

	cache.AddBlob(id, make([]byte, 10))

	if tx := dp.loki[toKey(id)]; len(tx) != 1 || tx[0] != Push{
		t.Fatal("Cache did not use policy properly")
	}

	if !found(){
		t.Fatal("Not found after adding to cache")
	}

	
	if tx := dp.loki[toKey(id)]; len(tx) != 2 || tx[1] != Visit{
		t.Fatal("Cache did not use policy properly")
	}
}

func TestCacheExit(t *testing.T){
	var(
		do = &DummyOperation{}
		id1 = []ops.Operation{do}
		id2 = append(id1, do)
		id3 = append(id2, do)
	)
	dp, cache := setup()

	cache.AddBlob(id1, make([]byte, 50))
	cache.AddBlob(id2, make([]byte, 40))
	cache.AddBlob(id3, make([]byte, 20))

	if dp.pops != 1{
		t.Fatal("Wrong amount of blobs removed from cache")
	}
}
