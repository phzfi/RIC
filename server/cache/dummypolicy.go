package cache

type DummyPolicy struct{}

func (q DummyPolicy) Visit(_ ImageInfo) {}

func (q DummyPolicy) Push(_ ImageInfo) {}

func (q DummyPolicy) Pop() (info ImageInfo) {
	info = ImageInfo{"", 0, 0, false}
	return
}
