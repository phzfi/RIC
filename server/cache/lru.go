package cache

import "sync"


var mutex = &sync.Mutex{}

type LRU struct {

	toList     map[string]*list
	head, tail list
}

func NewLRU() *LRU {

	lru := new(LRU)

	lru.toList = make(map[string]*list)
	lru.head.next = &lru.tail
	lru.tail.prev = &lru.head

	return lru
}

func (lru *LRU) Push(id string) {

	if _, ok := lru.toList[id]; ok {
		return
	}

	lru.push(id)
}

func (lru *LRU) push(id string) {
	l := list{id: id}

	prev := lru.last()
	l.prev = prev
	prev.next = &l

	l.next = &lru.tail
	lru.tail.prev = &l

	lru.toList[id] = &l
}

func (lru *LRU) Visit(id string) {

	mutex.Lock()
	defer mutex.Unlock()

	lru.toList[id].remove()
	lru.push(id)
}

func (lru *LRU) Pop() (id string) {

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
	id         string
}

func (l list) remove() {
	l.prev.next = l.next
	l.next.prev = l.prev
}
