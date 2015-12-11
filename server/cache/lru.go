package cache

func NewLRU(mm uint64) *Cache {
	return New(NewLRUPolicy(), mm)
}

type LRU struct {
	toList     map[ImageInfo]*list
	head, tail list
}

func NewLRUPolicy() *LRU {

	lru := new(LRU)

	lru.toList = make(map[ImageInfo]*list)
	lru.head.next = &lru.tail
	lru.tail.prev = &lru.head

	return lru
}

func (lru *LRU) Push(id ImageInfo) {

	l := list{id: id}

	next := &lru.tail
	l.next = next
	next.prev = &l

	prev := lru.tail.prev
	l.prev = prev
	prev.next = &l

	lru.toList[id] = &l
}

func (lru *LRU) Visit(id ImageInfo) {
	lru.toList[id].remove()
	lru.Push(id)
}

func (lru *LRU) Pop() (id ImageInfo) {

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
	id         ImageInfo
}

func (l list) remove() {

	if l.prev != nil {
		l.prev.next = l.next
	}
	if l.next != nil {
		l.next.prev = l.prev
	}
}
