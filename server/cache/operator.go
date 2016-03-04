package cache

import "github.com/phzfi/RIC/server/ops"
import "github.com/phzfi/RIC/server/images"

type Operator struct {
	cache  Cacher
	tokens TokenPool
}

type Cacher interface {
	GetBlob([]ops.Operation) (images.ImageBlob, bool)
	AddBlob([]ops.Operation, images.ImageBlob)
}

func MakeOperator(mm uint64) Operator {
	return Operator{
		HybridCache{
			NewLRU(mm),
			NewDiskCache("/tmp/RIC_diskcache", 1024*1024*1024*4),
		},
		MakeTokenPool(2),
	}
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
		o.tokens.Borrow()
		defer o.tokens.Return()

		//Check if some other thread already cached the image while we were blocked
		if blob, found := o.cache.GetBlob(operations); found {
			return blob, nil
		}

		img := images.NewImage()
		defer img.Destroy()

		if start != 0 {
			img.FromBlob(startimage)
		}

		o.applyOpsToImage(operations[start:], img)
		blob = img.Blob()

		o.cache.AddBlob(operations, blob)
	}

	return
}

func (o Operator) applyOpsToImage(operations []ops.Operation, img images.Image) (err error) {
	for _, op := range operations {
		err = op.Apply(img)
		if err != nil {
			return
		}
	}
	return
}
