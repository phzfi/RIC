package images

import (
	"bytes"
	"github.com/valyala/fasthttp"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func stopServer(ln net.Listener) {
	ln.Close()
	time.Sleep(100 * time.Millisecond)
}

func HandleTest(ctx *fasthttp.RequestCtx) {
	reader, err := os.Open("../testimages/loadimage/test.jpg")
	if err != nil {
		return
	}
	// Remember to free resources after you're done
	defer reader.Close()

	buffer := bytes.NewBuffer([]byte{})
	buffer.ReadFrom(reader)
	blob := buffer.Bytes()

	ctx.Write(blob)
}

func status404(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.WriteString("Image not found!")
}

func TestImageWeb(t *testing.T) {
	server := fasthttp.Server{
		Handler: HandleTest,
	}
	ln, _ := net.Listen("tcp", ":8009")
	go server.Serve(ln)
	defer stopServer(ln)
	time.Sleep(100 * time.Millisecond)

	img := NewImage()
	defer img.Destroy()
	err := img.FromWeb("http://localhost:8009/mikäliekuva.jpg")
	if err != nil {
		t.Fatal(err)
	}

	img_cmp := NewImage()
	defer img_cmp.Destroy()
	err = img_cmp.FromFile(filepath.FromSlash("../testimages/loadimage/test.jpg"))
	if err != nil {
		t.Fatal(err)
	}

	blob := img.Blob()
	blob_cmp := img_cmp.Blob()

	if len(blob) != len(blob_cmp) {
		t.Fatal("Image size different")
	}
	for i, v := range blob {
		if blob_cmp[i] != v {
			t.Fatal("Different image")
		}
	}

}

func TestImageWebWrongImage(t *testing.T) {
	server := fasthttp.Server{
		Handler: HandleTest,
	}
	ln, _ := net.Listen("tcp", ":8009")
	go server.Serve(ln)
	defer stopServer(ln)
	time.Sleep(100 * time.Millisecond)

	img := NewImage()
	defer img.Destroy()
	err := img.FromWeb("http://localhost:8009/mikäliekuva.jpg")
	if err != nil {
		t.Fatal(err)
	}

	img_cmp := NewImage()
	defer img_cmp.Destroy()
	err = img_cmp.FromFile(filepath.FromSlash("../testimages/loadimage/test.png"))
	if err != nil {
		t.Fatal(err)
	}

	blob := img.Blob()
	blob_cmp := img_cmp.Blob()

	if len(blob) != len(blob_cmp) {
		return
	}
	for i, v := range blob {
		if blob_cmp[i] != v {
			return
		}
	}

	t.Fatal("Images are same")

}

func TestImageWeb404(t *testing.T) {
	server := fasthttp.Server{
		Handler: status404,
	}
	ln, _ := net.Listen("tcp", ":8006")
	go server.Serve(ln)
	defer stopServer(ln)
	time.Sleep(100 * time.Millisecond)

	img := NewImage()
	defer img.Destroy()
	err := img.FromWeb("http://localhost:8006/mikäliekuva.jpg")
	if err == nil {
		t.Fatal("LoadImageWeb didn't return error when 404 received")
	}

}
