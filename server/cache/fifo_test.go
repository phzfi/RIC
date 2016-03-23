package cache

import "testing"

func TestFIFOUnderflow(t *testing.T) {
	fifo := FIFO{}

	defer func() {
		if recover() == nil {
			t.Fatal("No panic.")
		}
	}()

	fifo.Pop()
}

func TestFIFO(t *testing.T) {
	fifo := FIFO{}
	data := []string{"a", "b", "c", "d", "e", "f", "g"}

	add := func(data []string) {
		for _, d := range data {
			fifo.Push(d)
		}
	}

	remove := func(data []string) {
		for _, d := range data {
			if p := fifo.Pop(); p != d {
				t.Fatalf("Expected %s, got %s.", d, p)
			}
		}
	}

	half := len(data) / 2

	add(data[:half])
	remove(data[:half-1])
	add(data[half:])
	remove(data[half-1:])
}
