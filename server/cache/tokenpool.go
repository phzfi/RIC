package cache

type token struct{}

type TokenPool chan token

func MakeTokenPool(size int) (t TokenPool) {
	t = make(TokenPool, size)

	for i := 0; i < 2; i++ {
		t <- token{}
	}

	return
}

func (t TokenPool) Borrow() {
	<-t
}

func (t TokenPool) Return() {
	t <- token{}
}
