package cache

import "github.com/phzfi/RIC/server/ops"
import "github.com/phzfi/RIC/server/images"

type Operator struct {
	cache  *Cache
	tokens chan bool
}

func MakeOperator(mm uint64) Operator {
	o := Operator{NewLRU(mm), make(chan bool, 3)}
	for i := 0; i < 2; i++ {
		o.tokens <- true
	}
	return o
}

func (o Operator) GetBlob(operations ...ops.Operation) (blob images.ImageBlob, err error) {

	blob, found := o.cache.GetBlob(operations)
	if found {
		return blob, nil
	}

	t := <- o.tokens
	img := images.NewImage()
	defer func(){
		img.Destroy()
		o.tokens <- t
	}()

	for _, op := range operations {
		err = op.Apply(img)
		if err != nil {
			return
		}
	}

	blob = img.Blob()
	o.cache.AddBlob(operations, blob)
	return

}
