package operator

import (
	"errors"
	"testing"

	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
)

type BrokenOp struct{}

func (BrokenOp) Marshal() string { return "broken" }
func (BrokenOp) Apply(img images.Image) error { return errors.New("broken") }

// TestMakeBlob_Error tests MakeBlob with an operation that returns error
func TestMakeBlob_Error(t *testing.T) {
	op := MakeDefault(1000, "/tmp/operator_test_error", 1)

	_, err := op.GetBlob(BrokenOp{})
	if err == nil {
		t.Error("MakeBlob should return error for broken operation")
	}
}

// TestGetBlob_AlreadyCached tests getting a blob that's already cached
func TestGetBlob_AlreadyCached(t *testing.T) {
	op := MakeDefault(1000, "/tmp/operator_test_cached", 1)
	source := ops.MakeImageSource()
	source.AddRoot("/app/server/testimages/server/")
	loadOp := source.LoadImageOp("01.jpg")

	blob1, err := op.GetBlob(loadOp, ops.Resize{Width: 100, Height: 100})
	if err != nil {
		t.Fatalf("First GetBlob failed: %v", err)
	}

	// Second call should get from cache
	blob2, err := op.GetBlob(loadOp, ops.Resize{Width: 100, Height: 100})
	if err != nil {
		t.Fatalf("Second GetBlob failed: %v", err)
	}

	if len(blob1) != len(blob2) {
		t.Error("Cached blob should be identical")
	}
}

// TestGetBlob_PartiallyCached tests when some operations are cached
func TestGetBlob_PartiallyCached(t *testing.T) {
	op := MakeDefault(1000, "/tmp/operator_test_partial", 1)
	source := ops.MakeImageSource()
	source.AddRoot("/app/server/testimages/server/")
	loadOp := source.LoadImageOp("01.jpg")

	// First call with load + resize
	_, err := op.GetBlob(loadOp, ops.Resize{100, 100})
	if err != nil {
		t.Fatalf("First GetBlob failed: %v", err)
	}

	// Second call with load + resize + crop (partially cached: load+resize is cached)
	_, err = op.GetBlob(loadOp, ops.Resize{Width: 100, Height: 100}, ops.Crop{Width: 50, Height: 50, X: 0, Y: 0})
	if err != nil {
		t.Fatalf("Second GetBlob failed: %v", err)
	}
}
