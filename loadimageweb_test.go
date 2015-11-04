package main 

import( 
   "net/http"
   "testing"
   "fmt"
   "path/filepath"
   "gopkg.in/tylerb/graceful.v1"
   "time"
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
    fmt.Println("Ready to serve...")
    go func() {
        server.ListenAndServe()
    }()
    
    time.Sleep(time.Second)
    server.Stop(time.Second)

    fmt.Println("Server is up and running.")
    //if(err != nil) {
    //        t.Fatal(err)
    //}


}

