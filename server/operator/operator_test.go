package operator

import (
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/testutils"
	"testing"
)

const cacheFolder = "/tmp/operatortests"

func TestAlreadyCached(t *testing.T) {
	var log, log2 []int

	operator := MakeDefault(1000, cacheFolder)

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

func TestOperator(t *testing.T) {
	var log []int
	operations := []ops.Operation{
		&DummyOperation{&log, 0},
		&DummyOperation{&log, 1},
		&DummyOperation{&log, 2},
	}

	testutils.RemoveContents(cacheFolder)
	operator := MakeDefault(512*1024*1024, cacheFolder)

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
	testutils.RemoveContents(cacheFolder)

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
	operator := MakeDefault(512*1024*1024, cacheFolder)

	// Channel to track amount of completed operations
	c := make(chan bool, len(operations))

	// Launch operations simultaneously
	for _, ops := range operations {
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
