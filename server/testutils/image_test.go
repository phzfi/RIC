package testutils

import (
	"errors"
	"os"
	"testing"

	"github.com/phzfi/RIC/server/images"
)

func TestCheckDistortion(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	err = CheckDistortion(blob, "/app/server/testimages/server/01.jpg", 0.1, "/tmp/test_distortion.jpg")
	if err != nil {
		t.Errorf("CheckDistortion failed for identical images: %v", err)
	}
}

func TestCheckImage(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCase{
		Testfn: "/app/server/testimages/server/01.jpg",
		Reffn:  "/app/server/testimages/server/01.jpg",
		Resfn:   "/tmp/test_check_image.jpg",
	}

	err = CheckImage(blob, testCase, 0.1, CheckFormatFunc(FormatTestCase{testCase, "JPEG"}))
	if err != nil {
		t.Errorf("CheckImage failed: %v", err)
	}
}

func TestCheckFormatFunc(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCase{
		Testfn: "/app/server/testimages/server/01.jpg",
		Reffn:  "/app/server/testimages/server/01.jpg",
		Resfn:   "/tmp/test_format.jpg",
	}

	c := FormatTestCase{testCase, "JPEG"}

	err = CheckImage(blob, testCase, 0.1, CheckFormatFunc(c))
	if err != nil {
		t.Errorf("CheckFormatFunc failed: %v", err)
	}

	cWrong := FormatTestCase{testCase, "PNG"}
	err = CheckImage(blob, testCase, 0.1, CheckFormatFunc(cWrong))
	if err == nil {
		t.Error("CheckFormatFunc should fail with wrong format")
	}
}

func TestCheckSizeFunc(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCase{
		Testfn: "/app/server/testimages/server/01.jpg",
		Reffn:  "/app/server/testimages/server/01.jpg",
		Resfn:   "/tmp/test_size.jpg",
	}

	c := SizeTestCase{testCase, 900, 1200}

	err = CheckImage(blob, testCase, 0.1, CheckSizeFunc(c))
	if err != nil {
		t.Errorf("CheckSizeFunc failed: %v", err)
	}

	cWrong := SizeTestCase{testCase, 100, 100}
	err = CheckImage(blob, testCase, 0.1, CheckSizeFunc(cWrong))
	if err == nil {
		t.Error("CheckSizeFunc should fail with wrong size")
	}
}

func TestCheckAllFunc(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCaseAll{
		TestCase: TestCase{
			Testfn: "/app/server/testimages/server/01.jpg",
			Reffn:  "/app/server/testimages/server/01.jpg",
			Resfn:   "/tmp/test_all.jpg",
		},
		Format: "JPEG",
		W:      900,
		H:      1200,
	}

	err = CheckImage(blob, testCase.TestCase, 0.1, CheckAllFunc(testCase))
	if err != nil {
		t.Errorf("CheckAllFunc failed: %v", err)
	}
}

func TestFormatTest(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCase{
		Testfn: "/app/server/testimages/server/01.jpg",
		Reffn:  "/app/server/testimages/server/01.jpg",
		Resfn:   "/tmp/test_format_test.jpg",
	}

	c := FormatTestCase{testCase, "JPEG"}
	err = FormatTest(c, blob, 0.1)
	if err != nil {
		t.Errorf("FormatTest failed: %v", err)
	}
}

func TestSizeTest(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCase{
		Testfn: "/app/server/testimages/server/01.jpg",
		Reffn:  "/app/server/testimages/server/01.jpg",
		Resfn:   "/tmp/test_size_test.jpg",
	}

	c := SizeTestCase{testCase, 900, 1200}
	err = SizeTest(c, blob, 0.1)
	if err != nil {
		t.Errorf("SizeTest failed: %v", err)
	}
}

func TestTestAll(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCaseAll{
		TestCase: TestCase{
			Testfn: "/app/server/testimages/server/01.jpg",
			Reffn:  "/app/server/testimages/server/01.jpg",
			Resfn:   "/tmp/test_all_test.jpg",
		},
		Format: "JPEG",
		W:      900,
		H:      1200,
	}

	err = TestAll(testCase, blob, 0.1)
	if err != nil {
		t.Errorf("TestAll failed: %v", err)
	}
}

func TestCheckDistortion_InvalidRef(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	err = CheckDistortion(blob, "/nonexistent/path.jpg", 0.1, "/tmp/test_distortion_invalid.jpg")
	if err == nil {
		t.Error("CheckDistortion should fail with invalid ref path")
	}
}

func TestCheckDistortion_InvalidBlob(t *testing.T) {
	err := CheckDistortion([]byte("invalid blob"), "/app/server/testimages/server/01.jpg", 0.1, "/tmp/test_distortion_invalid_blob.jpg")
	if err == nil {
		t.Error("CheckDistortion should fail with invalid blob")
	}
}

func TestCheckImage_ErrorFromFunc(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCase{
		Testfn: "/app/server/testimages/server/01.jpg",
		Reffn:  "/app/server/testimages/server/01.jpg",
		Resfn:   "/tmp/test_check_image_error.jpg",
	}

	errorFunc := func(img images.Image) error {
		return errors.New("simulated error")
	}

	err = CheckImage(blob, testCase, 0.1, errorFunc)
	if err == nil {
		t.Error("CheckImage should fail when func returns error")
	}
}

func TestCheckAllFunc_WrongSize(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCaseAll{
		TestCase: TestCase{
			Testfn: "/app/server/testimages/server/01.jpg",
			Reffn:  "/app/server/testimages/server/01.jpg",
			Resfn:   "/tmp/test_all_wrong_size.jpg",
		},
		Format: "JPEG",
		W:      100,
		H:      100,
	}

	err = CheckImage(blob, testCase.TestCase, 0.1, CheckAllFunc(testCase))
	if err == nil {
		t.Error("CheckAllFunc should fail with wrong size")
	}
}

func TestCheckAllFunc_WrongFormat(t *testing.T) {
	img := images.NewImage()
	defer img.Destroy()
	err := img.FromFile("/app/server/testimages/server/01.jpg")
	if err != nil {
		t.Fatalf("Failed to load test image: %v", err)
	}
	blob := img.Blob()

	testCase := TestCaseAll{
		TestCase: TestCase{
			Testfn: "/app/server/testimages/server/01.jpg",
			Reffn:  "/app/server/testimages/server/01.jpg",
			Resfn:   "/tmp/test_all_wrong_format.jpg",
		},
		Format: "PNG",
		W:      900,
		H:      1200,
	}

	err = CheckImage(blob, testCase.TestCase, 0.1, CheckAllFunc(testCase))
	if err == nil {
		t.Error("CheckAllFunc should fail with wrong format")
	}
}

func TestRemoveContents_Errors(t *testing.T) {
	// Test with non-existent directory
	err := RemoveContents("/nonexistent/dir/path")
	if err == nil {
		t.Error("RemoveContents should fail with non-existent directory")
	}

	// Test with a file instead of directory
	err = RemoveContents("/app/server/testimages/server/01.jpg")
	if err == nil {
		t.Error("RemoveContents should fail when given a file path")
	}

	// Test with directory that can't be read
	// This is difficult to test directly since most dirs are readable
	// We'll create a dir and then make it unreadable if possible
	testDir := "/tmp/testutils_test_rm_errors"
	os.MkdirAll(testDir, 0755)
	
	// Try to remove contents from a file (not dir)
	err = RemoveContents("/app/server/testimages/server/01.jpg")
	if err == nil {
		t.Error("RemoveContents should fail with file path")
	}
}
