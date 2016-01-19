package main

import (
	"flag"
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/logging"
	"github.com/valyala/fasthttp"
	"gopkg.in/tylerb/graceful.v1"
	//"io"
	"log"
	//"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

// MyHandler type is used to encompass HandlerFunc interface.
// In the future this type will probably contain pointers to
// services provided by this program (image cache).
type MyHandler struct {

	// Service started
	started time.Time

	// Request count (statistics)
	requests uint64

	images cache.AmbiguousSizeImageCache
}

// ServeHTTP is called whenever there is a new request.
// This is quite similar to JavaEE Servlet interface.
// TODO: Check that ServeHTTP is called inside a goroutine?
func (h *MyHandler) ServeHTTP(/*writer http.ResponseWriter, request *http.Request*/ctx *fasthttp.RequestCtx) {
	//method := request.Method

	// In the future we can use requester can detect request spammers!
	// requester := request.RemoteAddr

	// Increase request count
	count := &(h.requests)
	atomic.AddUint64(count, 1)

	if /*method == "GET"*/ctx.IsGet() {

		//url := request.URL
		//filename := url.Path
		url := ctx.URI()
		filename := string(ctx.Path())

		// GET parameters
		//query := url.Query()
		query := url.QueryArgs()

		qw, _ := query.GetUint("width")
		qh, _ := query.GetUint("height")
		qm, _ := query.GetUint("mode")
		uqw := uint(qw)
		uqh := uint(qh)
		uqm := uint(qm)
		puqw := &uqw
		puqh := &uqh
		puqm := &uqm
		h.RetrieveImage(
			//writer, filename,
			//getUintParam(query, "width"),
			//getUintParam(query, "height"),
			//getStringParam(query, "mode"))
			ctx, filename,
			puqw,
			puqh,
			puqm)

	} else if ctx.IsPost()/*method == "POST"*/ {
		// POST is currently unused so we can use this for testing
		h.RetrieveHello(ctx)
	}
}

// Returns a request parameter as *uint; nil if the parameter is not properly specified.
func getUintParam(params map[string][]string, name string) (result *uint) {

	if values := params[name]; len(values) != 0 {
		asUint, err := strconv.ParseUint(values[0], 10, 32)
		if err == nil {
			u := uint(asUint)
			result = &u
		}
	}
	return
}

// Returns a request parameter as *string; nil if the parameter is not properly specified.
func getStringParam(params map[string][]string, name string) (result *string) {
	if values := params[name]; len(values) != 0 {
		result = &values[0]
	}
	return
}

// Respond to POST message by saying Hello
func (h MyHandler) RetrieveHello(/*writer http.ResponseWriter*/ctx *fasthttp.RequestCtx) {
	_, err := ctx.WriteString("Hello world!")
	/*
	_, err := io.WriteString(writer, result)
	*/
	if err != nil {
		log.Println(err)
	}
}

// Write image by filename into ResponseWriter with the
// desired width and height being pointed to. If there
// are no desired width or height, that parameter is nil.
func (h *MyHandler) RetrieveImage(/*writer http.ResponseWriter*/ctx *fasthttp.RequestCtx,
	filename string,
	width *uint,
	height *uint,
	mode *string) {

	// TODO: filename must not be interpret as "absolute"
	// implement a type that will abstract away the filesystem.
	logging.Debug("Find: " + filename)

	// Get cache
	bank := h.images

	// Load the image
	blob, err := bank.GetImage(filename, width, height, mode)
	if err != nil {
		// TODO:
		// Classify different possible errors more but make sure
		// no "internal" information is leaked.
		//writer.WriteHeader(http.StatusNotFound)
		//io.WriteString(writer, "Image not found!")
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.WriteString("Image not found!")
		logging.Debug(err)
		return
	}
	//writer.Write(blob)
	ctx.Write(blob)
}

// Create a new graceful server and configure it.
// This does not run the server however.
func NewServer(maxMemory uint64) (*graceful.Server, *MyHandler) {

	cacher := cache.AmbiguousSizeImageCache{cache.NewLRU(maxMemory)}

	// Add roots
	// TODO: This must be externalized outside the source code.
	if cacher.AddRoot("/var/www") != nil {
		log.Fatal("Root not added /var/www")
	}

	if cacher.AddRoot(".") != nil {
		log.Println("Root not added .")
	}

	// Configure handler
	handler := &MyHandler{
		requests: 0,
		images:   cacher,
	}

	// Configure server
	server := &graceful.Server{
		Timeout: 8 * time.Second,
		Server: &fasthttp.Server{
			Addr:           ":8005",
			Handler:        handler,
			ReadTimeout:    8 * time.Second,
			WriteTimeout:   8 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
	return server, handler
}

func main() {

	// CLI arguments
	mem := flag.Uint64("m", 500*1024*1024, "Sets the maximum memory to be used for caching images in bytes. Does not account for memory consumption of other things.")
	flag.Parse()

	server, handler := NewServer(*mem)

	log.Println("Server starting...")
	handler.started = time.Now()
	err := server.ListenAndServe()
	end := time.Now()

	// Get number of requests
	requests := strconv.FormatUint((*handler).requests, 10)

	// Calculate the elapsed time
	duration := end.Sub(handler.started)
	log.Println("Server requests: " + requests)
	log.Println("Server uptime: " + duration.String())

	// Log errors
	if err != nil {
		log.Fatal(err)
	}
}
