package ops

import (
	"testing"
	"time"
	"github.com/phzfi/RIC/server/images"
	"fmt"
)


func TestPingRoots(t *testing.T) {

	is := MakeImageSource()
	err := is.AddRoot("../")
	if err != nil {
		t.Fatal(err)
	}

	img := images.NewImage()
	defer img.Destroy()
	t0 := time.Now()
	err = is.searchRoots("testimages/loadimage/test.jpg", img)
	if err != nil {
		t.Fatal(err)
	}
	t1 := time.Now().Sub(t0)
	
	img = images.NewImage()
	defer img.Destroy()
	t0 = time.Now()
	err = is.pingRoots("testimages/loadimage/test.jpg", img)
	if err != nil {
		t.Fatal(err)
	}
	t2 := time.Now().Sub(t0)
	
	img = images.NewImage()
	defer img.Destroy()
	t0 = time.Now()
	err = is.searchRoots("testimages/loadimage/test.jpg", img)
	if err != nil {
		t.Fatal(err)
	}
	t3 := time.Now().Sub(t0)
	
	if t1 < t2 {
		t.Fatal(fmt.Sprintf("pingRoots is slower than searchRoots! pingRoots: %v, searchRoots: %v", t2, t1))
	}
	
	if t3 < t2 {
		t.Fatal(fmt.Sprintf("pingRoots is slower than searchRoots! pingRoots: %v, searchRoots: %v", t2, t3))
	}
	return
}
