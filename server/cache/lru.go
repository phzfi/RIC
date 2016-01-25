package cache

import "sync"

func NewLRU(mm uint64) *Cache {
	return NewCache(NewLRUPolicy(), mm)
}

type LRU struct {
	sync.Mutex

	toList     map[cacheKey]*list
	head, tail list
}

func NewLRUPolicy() *LRU {

	lru := new(LRU)

	lru.toList = make(map[cacheKey]*list)
	lru.head.next = &lru.tail
	lru.tail.prev = &lru.head

	return lru
}

func (lru *LRU) Push(id cacheKey) {

	l := list{id: id}

	prev := lru.tail.prev
	l.prev = prev
	prev.next = &l

	l.next = &lru.tail
	lru.tail.prev = &l

	lru.toList[id] = &l
}

func (lru *LRU) Visit(id cacheKey) {

	lru.Lock()
	defer lru.Unlock()

	lru.toList[id].remove()
	lru.Push(id)
}

func (lru *LRU) Pop() (id cacheKey) {

	first := lru.first()
	if first == &lru.tail {
		panic("LRU underflow")
	}

	id = first.id
	first.remove()
	delete(lru.toList, id)

	return
}

func (lru LRU) first() *list {
	return lru.head.next
}

func (lru LRU) last() *list {
	return lru.tail.prev
}

type list struct {
	next, prev *list
	id         cacheKey
}

func (l list) remove() {
	l.prev.next = l.next
	l.next.prev = l.prev
}
