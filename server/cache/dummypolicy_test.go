package cache

type Log map[string][]uint

const (
	Visit = iota
	Push
	Pop
)

type DummyPolicy struct {
	fifo Policy

	loki Log
	pops int
}

func (d DummyPolicy) Visit(k string) {
	d.log(k, Visit)
	d.fifo.Visit(k)
}

func (d DummyPolicy) log(k string, t uint) {
	d.loki[k] = append(d.loki[k], t)
}

func (d DummyPolicy) Push(k string) {
	d.log(k, Push)
	d.fifo.Push(k)
}

func (d *DummyPolicy) Pop() string {
	d.pops += 1
	return d.fifo.Pop()
}

func NewDummyPolicy() *DummyPolicy {
	return &DummyPolicy{fifo: &FIFO{}, loki: make(Log)}
}
