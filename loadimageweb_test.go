package main 

import( 
   "net/http"
   "testing"
   "path/filepath"
   "gopkg.in/tylerb/graceful.v1"
)


type TestHandler struct {
}

func(h *TestHandler) ServeHTTP(w http.ResponseWriter, request *http.Request) {
    image, err := LoadImage(filepath.FromSlash("testimages/loadimage/test.jpg"))
    if(err != nil){
        return
    }
    w.Write(image)
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
    
    imageblob, err := LoadImageWeb("http://localhost:8006/mikäliekuva.jpg")
    if err != nil {
        t.Fatal(err)
    }

    imageblob_cmp, err := LoadImage(filepath.FromSlash("testimages/loadimage/test.jpg"))
    if err != nil {
        t.Fatal(err)
    }

    if len(imageblob) != len(imageblob_cmp) {
        t.Fatal("Image size different")
    }
    for i, v := range imageblob {
        if imageblob_cmp[i] != v {
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
    
    imageblob, err := LoadImageWeb("http://localhost:8006/mikäliekuva.jpg")
    if err != nil {
        t.Fatal(err)
    }

    imageblob_cmp, err := LoadImage(filepath.FromSlash("testimages/loadimage/test.png"))
    if err != nil {
        t.Fatal(err)
    }

    if len(imageblob) != len(imageblob_cmp) {
        return
    }
    for i, v := range imageblob {
        if imageblob_cmp[i] != v {
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
