package main

import (
	"fmt"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/testutils"
	"testing"
)

type CommonTestCase struct {
	test testutils.TestCaseAll
	op   ops.Operation
}

type CommonTest func(CommonTestCase) error

func threadedTesting(cases []CommonTestCase, test CommonTest) (err error) {
	sem := make(chan error, len(cases))
	for _, c := range cases {
		go func(tc CommonTestCase) {
			sem <- test(tc)
		}(c)
	}
	for range cases {
		var verr = <-sem
		if verr != nil && err == nil {
			// Pick the first error but wait for termination
			err = verr
		}
	}
	return
}

func TestOperatorConvert(t *testing.T) {

	operator, src := SetupOperatorSource()


	testfolder := "/ric/assets/test_assets/testimages/convert/"
	testimage := testfolder + "toconvert.jpg"
	testimage2 := testfolder + "toconvert2.jpg"
	testimage3 := testfolder + "toconvert3.jpg"
	testimage4 := testfolder + "toconvert4.jpg"
	resfolder := "/ric/assets/test_assets/testresults/common/"
	tolerance := 0.002

	var conv = func(a, b, c, d string) CommonTestCase {
		va := testutils.TestCase{a, "", testfolder + b, resfolder + c}
		vb := testutils.TestCaseAll{va, d, -1, -1}
		return CommonTestCase{vb, ops.Convert{d}}
	}

	cases := []CommonTestCase{
		conv(testimage, "converted.jpg", "converted.jpg", "JPEG"),
		//conv(testimage, "converted.webp", "converted.webp", "WEBP"),
		conv(testimage, "converted.tiff", "converted.tiff", "TIFF"),
		conv(testimage, "converted.gif", "converted.gif", "GIF"),
		conv(testimage, "converted.png", "converted.png", "PNG"),
		conv(testimage, "converted.bmp", "converted.bmp", "BMP"),
		conv(testimage2, "converted2.jpg", "converted2.jpg", "JPEG"),
		//conv(testimage2, "converted2.webp", "converted2.webp", "WEBP"),
		conv(testimage3, "converted3.jpg", "converted3.jpg", "JPEG"),
		//conv(testimage3, "converted3.webp", "converted3.webp", "WEBP"),
		conv(testimage4, "converted4.jpg", "converted4.jpg", "JPEG"),
		//conv(testimage4, "converted4.webp", "converted4.webp", "WEBP"),
	}

	var test = func(c CommonTestCase) (err error) {

		var vt = c.test
		var vo = c.op.(ops.Convert)
		logging.Debug(fmt.Sprintf("Testing convert: %v, %v, %v, %v", vt.TestFilename, vt.ReferenceFilename, vt.Format, vt.ResultFilename))

		blob, err := operator.GetBlob("namespace", src.LoadImageOp(vt.TestFilename), vo)
		if err != nil {
			return
		}

		var ft = testutils.FormatTestCase{testutils.TestCase{vt.TestFilename, "", vt.ReferenceFilename, vt.ResultFilename}, vt.Format}
		err = testutils.FormatTest(ft, blob, tolerance)
		return
	}

	var verr = threadedTesting(cases, test)
	if verr != nil {
		t.Fatal(verr)
	}
}

func TestOperatorResize(t *testing.T) {

	operator, src := SetupOperatorSource()

	testfolder := "/ric/assets/test_assets/testimages/resize/"
	testimage := testfolder + "toresize.jpg"
	testimage2 := testfolder + "toresize2.jpg"
	testimage3 := testfolder + "toresize3.jpg"
	testimage4 := testfolder + "toresize4.jpg"
	resfolder := "/ric/assets/test_assets/testresults/common/"
	tolerance := 0.002

	var res = func(a, b, c string, d, e int) CommonTestCase {
		va := testutils.TestCase{a, "", testfolder + b, resfolder + c}
		vb := testutils.TestCaseAll{va, "Whatever", d, e}
		return CommonTestCase{vb, ops.Resize{d, e}}
	}

	cases := []CommonTestCase{
		res(testimage, "1_100x100.jpg", "1_100x100.jpg", 100, 100),
		res(testimage, "1_200x200.jpg", "1_200x200.jpg", 200, 200),
		res(testimage, "1_300x400.jpg", "1_300x400.jpg", 300, 400),
		res(testimage, "1_500x200.jpg", "1_500x200.jpg", 500, 200),
		res(testimage, "1_30x20.jpg", "1_30x20.jpg", 30, 20),
		res(testimage, "1_600x600.jpg", "1_600x600.jpg", 600, 600),
		res(testimage2, "2_100x100.jpg", "2_100x100.jpg", 100, 100),
		res(testimage2, "2_200x200.jpg", "2_200x200.jpg", 200, 200),
		res(testimage3, "3_100x100.jpg", "3_100x100.jpg", 100, 100),
		res(testimage3, "3_200x200.jpg", "3_200x200.jpg", 200, 200),
		res(testimage4, "4_100x100.jpg", "4_100x100.jpg", 100, 100),
		res(testimage4, "4_200x200.jpg", "4_200x200.jpg", 200, 200),
	}

	var test = func(c CommonTestCase) (err error) {
		var vt = c.test
		var vo = c.op.(ops.Resize)
		logging.Debug(fmt.Sprintf("Testing resize: %v, %v, %v, %v, %v", vt.TestFilename, vt.ReferenceFilename, vt.W, vt.H, vt.ResultFilename))

		blob, err := operator.GetBlob("namespace", src.LoadImageOp(vt.TestFilename), vo)
		if err != nil {
			return
		}

		var rt = testutils.SizeTestCase{testutils.TestCase{vt.TestFilename, "", vt.ReferenceFilename, vt.ResultFilename}, vt.W, vt.H}
		err = testutils.SizeTest(rt, blob, tolerance)
		return
	}

	var verr = threadedTesting(cases, test)
	if verr != nil {
		t.Fatal(verr)
	}
}

func TestOperatorLiquidRescale(t *testing.T) {

	operator, src := SetupOperatorSource()

	testfolder := "/ric/assets/test_assets/testimages/resize/"
	testimage := testfolder + "toresize.jpg"
	testimage2 := testfolder + "toresize2.jpg"
	resfolder := "/ric/assets/test_assets/testresults/common/"
	tolerance := 0.05

	var res = func(a, b, c string, d, e int) CommonTestCase {
		va := testutils.TestCase{a, "", testfolder + b, resfolder + c}
		vb := testutils.TestCaseAll{va, "Whatever", d, e}
		return CommonTestCase{vb, ops.LiquidRescale{d, e}}
	}

	cases := []CommonTestCase{
		res(testimage, "liquid1_100x100.jpg", "liquid1_100x100.jpg", 100, 100),
		res(testimage, "liquid1_500x200.jpg", "liquid1_500x200.jpg", 500, 200),
		res(testimage2, "liquid2_200x200.jpg", "liquid2_200x200.jpg", 200, 200),
	}

	var test = func(c CommonTestCase) (err error) {
		var vt = c.test
		var vo = c.op.(ops.LiquidRescale)
		logging.Debug(fmt.Sprintf("Testing resize: %v, %v, %v, %v, %v", vt.TestFilename, vt.ReferenceFilename, vt.W, vt.H, vt.ResultFilename))

		blob, err := operator.GetBlob("namespace", src.LoadImageOp(vt.TestFilename), vo)
		if err != nil {
			return
		}

		var rt = testutils.SizeTestCase{testutils.TestCase{vt.TestFilename, "", vt.ReferenceFilename, vt.ResultFilename}, vt.W, vt.H}
		err = testutils.SizeTest(rt, blob, tolerance)
		return
	}

	var verr = threadedTesting(cases, test)
	if verr != nil {
		t.Fatal(verr)
	}
}
