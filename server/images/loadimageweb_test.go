package images

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"github.com/valyala/fasthttp"
	"net"
	"time"
)

type TestHandler struct {}
type NotFoundHandler struct {}


func stopServer(ln net.Listener) {
	ln.Close()
	time.Sleep(100 * time.Millisecond)
}


func (h *TestHandler) ServeHTTP(ctx *fasthttp.RequestCtx) {
	reader, err := os.Open("../testimages/loadimage/test.jpg")
	if err != nil {
		return
	}
	// Remember to free resources after you're done
	defer reader.Close()

	buffer := bytes.NewBuffer([]byte{})
	buffer.ReadFrom(reader)
	blob := ImageBlob(buffer.Bytes())

	ctx.Write(blob)
}

func (h *NotFoundHandler) Serve404 (ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.WriteString("Image not found!")
}


func TestImageWeb(t *testing.T) {
	handler := &TestHandler{}
	server := &fasthttp.Server{
		Handler: handler.ServeHTTP,
	}
	ln, _ := net.Listen("tcp", ":8005")
	go server.Serve(ln)
	defer stopServer(ln)
	time.Sleep(100 * time.Millisecond)

	image, err := LoadImageWeb("http://localhost:8005/mikäliekuva.jpg")
	if err != nil {
		t.Fatal(err)
	}

	image_cmp, err := LoadImage(filepath.FromSlash("../testimages/loadimage/test.jpg"))
	if err != nil {
		t.Fatal(err)
	}

	blob := image.ToBlob()
	blob_cmp := image_cmp.ToBlob()

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
	handler := &TestHandler{}
	server := &fasthttp.Server{
		Handler: handler.ServeHTTP,
	}
	ln, _ := net.Listen("tcp", ":8005")
	go server.Serve(ln)
	defer stopServer(ln)
	time.Sleep(100 * time.Millisecond)


	image, err := LoadImageWeb("http://localhost:8005/mikäliekuva.jpg")
	if err != nil {
		t.Fatal(err)
	}

	image_cmp, err := LoadImage(filepath.FromSlash("../testimages/loadimage/test.png"))
	if err != nil {
		t.Fatal(err)
	}

	blob := image.ToBlob()
	blob_cmp := image_cmp.ToBlob()

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
	handler := &NotFoundHandler{}
	server := &fasthttp.Server{
		Handler: handler.Serve404,
	}
	ln, _ := net.Listen("tcp", ":8006")
	go server.Serve(ln)
	defer stopServer(ln)
	time.Sleep(100 * time.Millisecond)


	_, err := LoadImageWeb("http://localhost:8006/mikäliekuva.jpg")
	if err == nil {
		t.Fatal("LoadImageWeb didn't return error when 404 received")
	}

}
