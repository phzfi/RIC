package cache

import (
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"testing"
)

type DummyOperation struct {
	log  *[]int
	name int
}

func (o *DummyOperation) Apply(img images.Image) (images.Image, error) {
	*(o.log) = append(*(o.log), o.name)
	return img, nil
}

func TestOperator(t *testing.T) {
	var log []int
	operations := []ops.Operation{
		&DummyOperation{&log, 0},
		&DummyOperation{&log, 1},
		&DummyOperation{&log, 2},
	}
	operator := MakeOperator(512 * 1024 * 1024)

	_, err := operator.GetBlob(operations)
	if err != nil {
		t.Error(err)
	}

	if len(log) != 3 {
		t.Fatal("Too many or too few operations done")
	}
	for i, v := range log {
		if i != v {
			t.Fatal("Wrong operation")
		}
	}
}
