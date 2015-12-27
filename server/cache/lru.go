package cache

func NewLRU(mm uint64) ImageCache {
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

	prev := lru.tail.prev
	l.prev = prev
	prev.next = &l

	l.next = &lru.tail
	lru.tail.prev = &l

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
	l.prev.next = l.next
	l.next.prev = l.prev
}
