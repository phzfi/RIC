package operator

import (
	"errors"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/testutils"
	"testing"
)

const cacheFolder = "/tmp/operatortests"

var tokens = 3

func prepare() Operator {
	testutils.RemoveContents(cacheFolder)
	return MakeWithDefaultCacheSet(1000, cacheFolder, tokens)
}

func TestAlreadyCached(t *testing.T) {
	var log, log2 []int

	operator := prepare()

	operator.GetBlob(
		&DummyOperation{&log, 9},
		&DummyOperation{&log, 3},
	)
	operator.GetBlob(
		&DummyOperation{&log2, 9},
		&DummyOperation{&log2, 3},
	)

	if len(log2) != 0 {
		t.Fatal("Operator did not use a cached result, instead running the operations again.")
	}
}

func TestPartiallyCached(t *testing.T) {
	var log, log2 []int

	operator := prepare()

	operator.GetBlob(&DummyOperation{&log, 9})
	operator.GetBlob(
		&DummyOperation{&log2, 9},
		&DummyOperation{&log2, 3},
	)

	if len(log2) != 1 {
		t.Fatalf("Operator ran %d operations instead of 1.", len(log2))
	}
}

func TestOperator(t *testing.T) {
	var log []int
	operations := []ops.Operation{
		&DummyOperation{&log, 0},
		&DummyOperation{&log, 1},
		&DummyOperation{&log, 2},
	}

	operator := prepare()

	_, err := operator.GetBlob(operations...)
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

func TestDenyIdenticalOperations(t *testing.T) {
	var log []int

	// Many identical operations
	operations := [][]ops.Operation{
		{&DummyOperation{&log, 0}, &DummyOperation{&log, 0}},
		{&DummyOperation{&log, 0}, &DummyOperation{&log, 0}},
		{&DummyOperation{&log, 0}, &DummyOperation{&log, 0}},
		{&DummyOperation{&log, 0}, &DummyOperation{&log, 0}},
		{&DummyOperation{&log, 0}, &DummyOperation{&log, 0}},
		{&DummyOperation{&log, 0}, &DummyOperation{&log, 0}},
	}
	operator := prepare()

	// Channel to track amount of completed operations
	c := make(chan bool, len(operations))

	// Launch operations simultaneously
	for i := range operations {
		ops := operations[i]
		go func() {
			_, _ = operator.GetBlob(ops...)
			c <- true
		}()
	}

	// Wait for the operations to finish
	for i := 0; i < len(operations); i++ {
		<-c
	}

	// Only 2 operations should've been done - others found from cache
	if len(log) != 2 {
		t.Fatalf("%v operations done. Expected 2", len(log))
	}
}

type BrokenOperation struct{}

func (BrokenOperation) Marshal() string {
	return "broken"
}

func (BrokenOperation) Apply(image images.Image) error {
	return errors.New("This operation is broken")
}

func TestBrokenOperation(t *testing.T) {
	operator := prepare()
	_, err := operator.GetBlob(BrokenOperation{})
	if err == nil {
		t.Fatal("Broken operation did not return error")
	}
}
