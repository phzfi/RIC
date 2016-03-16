package cache

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"sync"
	"testing"
	"time"
)

type DummyOperation struct {
	log  *[]int
	name int
}

func (o *DummyOperation) Marshal() string {
	return fmt.Sprintf("test%v", o.name)
}

var logMutex *sync.Mutex = &sync.Mutex{}

func (o *DummyOperation) Apply(img images.Image) error {
	// Take some time for simult opers. tests
	time.Sleep(200 * time.Millisecond)
	logMutex.Lock()
	*(o.log) = append(*(o.log), o.name)
	logMutex.Unlock()
	return nil
}

func TestOperator(t *testing.T) {
	var log []int
	operations := []ops.Operation{
		&DummyOperation{&log, 0},
		&DummyOperation{&log, 1},
		&DummyOperation{&log, 2},
	}

	removeContents(cacheFolder)
	operator := MakeOperator(512*1024*1024, cacheFolder)

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
	operator := MakeOperator(512*1024*1024, cacheFolder)

	// Channel to track amount of completed operations
	c := make(chan bool, len(operations))

	// Launch operations simultaneously
	for i := 0; i < len(operations); i++ {
		opers := operations[i]
		go func() {
			_, _ = operator.GetBlob(opers...)
			c <- true
		}()
	}

	// Wait for the operations to finish
	for i := 0; i < len(operations); i++ {
		_ = <-c
	}

	// Only 2 operations should've been done - others found from cache
	if len(log) != 2 {
		t.Fatal(fmt.Sprintf("%v operations done. Expected 2", len(log)))
	}
}
