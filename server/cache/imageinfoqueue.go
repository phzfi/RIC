package cache

type imageInfoQueue struct {
	data       []imageInfo
	head, tail int
}

func (q *imageInfoQueue) Push(info imageInfo) {

	if size, arraySize := q.Len(), len(q.data); size + 2 > arraySize {

		new_data := make([]imageInfo, (size+1)*2)

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

func (q *imageInfoQueue) Pop() (info imageInfo) {

	if q.Len() == 0 {
		panic("Queue underflow")
	}

	info = q.data[q.head]
	q.head++
	q.head %= len(q.data)

	return
}

func (q imageInfoQueue) Len() int {
	size := q.tail - q.head
	if size < 0 {
		return len(q.data) + size
	}
	return size
}
