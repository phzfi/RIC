package cache

import "errors"

type FIFO struct {
	data       []string
	head, tail int
}

func (q *FIFO) Visit(_ string) {}

func (q *FIFO) Push(info string) {

	if size, arraySize := q.Len(), len(q.data); size+2 > arraySize {

		new_data := make([]string, (size+1)*2)

		for i := 0; i < size; i++ {
			new_data[i] = q.data[(q.head+i)%arraySize]
		}
		q.data = new_data
		q.head = 0
		q.tail = size
	}

	q.data[q.tail] = info
	q.tail++
	q.tail %= len(q.data)
}

func (q *FIFO) Pop() (info string, err error) {

	if q.Len() == 0 {
		err = errors.New("no items in queue")
		return
	}

	info = q.data[q.head]
	q.head++
	q.head %= len(q.data)

	return
}

func (q FIFO) Len() int {
	size := q.tail - q.head
	if size < 0 {
		return len(q.data) + size
	}
	return size
}
