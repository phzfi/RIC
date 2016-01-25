package cache


import "github.com/phzfi/RIC/server/ops"
import "github.com/phzfi/RIC/server/images"


type Operator struct{
	cache *Cache
}


func MakeOperator(mm uint64) Operator {
	return Operator{NewLRU(mm)}
}


func (o Operator) GetBlob(operations ...ops.Operation) (blob images.ImageBlob, err error) {
	
	blob, found := o.cache.GetBlob(operations)
	if found {
		return blob, nil
	}

	img := images.NewImage()
	defer img.Destroy()

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
