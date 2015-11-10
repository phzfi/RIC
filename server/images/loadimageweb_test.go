package images

import( 
   "net/http"
   "testing"
   "path/filepath"
   "gopkg.in/tylerb/graceful.v1"
   "os"
   "bytes"
)


type TestHandler struct {
}

func(h *TestHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
    reader, err := os.Open("../testimages/loadimage/test.jpg")
    if err != nil {
    	return
    }
    // Remember to free resources after you're done
    defer reader.Close()
    
    buffer := bytes.NewBuffer([]byte{})
    buffer.ReadFrom(reader)
    blob := ImageBlob(buffer.Bytes())

    w.Write(blob)
}


func TestImageWeb(t *testing.T) {
    handler := &TestHandler{
    }
    server := graceful.Server{
        Server: &http.Server{
            Addr: ":8006",
            Handler: handler,
        },
    }
    
    go server.ListenAndServe()
    defer server.Stop(0)
    
    image, err := LoadImageWeb("http://localhost:8006/mikäliekuva.jpg")
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
    handler := &TestHandler{
    }
    server := graceful.Server{
        Server: &http.Server{
            Addr: ":8006",
            Handler: handler,
        },
    }
    
    go server.ListenAndServe()
    defer server.Stop(0)
    
    image, err := LoadImageWeb("http://localhost:8006/mikäliekuva.jpg")
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

    server := graceful.Server{
        Server: &http.Server{
            Addr: ":8006",
            Handler: http.NotFoundHandler(),
        },
    }

    go server.ListenAndServe()
    defer server.Stop(0)

    _, err := LoadImageWeb("http://localhost:8006/mikäliekuva.jpg")
    if err == nil {
        t.Fatal("LoadImageWeb didn't return error when 404 received")
    }

}
