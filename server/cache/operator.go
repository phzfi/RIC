package cache

import "github.com/phzfi/RIC/server/ops"
import "github.com/phzfi/RIC/server/images"

type Operator struct {
	cache  *Cache
	tokens chan bool
}

func MakeOperator(mm uint64) Operator {
	o := Operator{NewLRU(mm), make(chan bool, 3)}
	// TODO: Currently only 2 simult. operations allowed. Increate tokens and make them configurable.
	for i := 0; i < 2; i++ {
		o.tokens <- true
	}
	return o
}

func (o Operator) GetBlob(operations ...ops.Operation) (blob images.ImageBlob, err error) {

	var startimage images.ImageBlob
	var start int

	for start = len(operations); start > 0; start-- {
		var found bool
		startimage, found = o.cache.GetBlob(operations[:start])
		if found {
			break
		}
	}

	if start == len(operations) {
		return startimage, nil
	} else {
		t := <-o.tokens
		defer func() { o.tokens <- t }()

		// Check if some other thread already cached the image while we were blocked
		if blob, found := o.cache.GetBlob(operations); found {
			return blob, nil
		}
		
		img := images.NewImage()
		defer img.Destroy()

		if start != 0 {
			img.FromBlob(startimage)
		}

		blob, err = o.applyOpsToImage(operations, start, img)
		if err != nil {
			return
		}
	}

	return
}

func (o Operator) applyOpsToImage(operations []ops.Operation, start int, img images.Image) (blob images.ImageBlob, err error) {
	for i, op := range operations[start:] {
		err = op.Apply(img)
		if err != nil {
			return
		}
		blob = img.Blob()
		o.cache.AddBlob(operations[:start + i + 1], blob)
	}
	return
}
