package operator

import (
	"fmt"
	"github.com/phzfi/RIC/server/images"
	"sync"
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
	time.Sleep(50 * time.Millisecond)
	logMutex.Lock()
	*(o.log) = append(*(o.log), o.name)
	logMutex.Unlock()
	return nil
}
