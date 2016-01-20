package main

import (
	"flag"
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/logging"
	"github.com/valyala/fasthttp"
	"log"
	"strconv"
	"sync/atomic"
	"time"
	"net"
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
func (h *MyHandler) ServeHTTP(ctx *fasthttp.RequestCtx) {

	// In the future we can use requester can detect request spammers!
	// requester := ctx.RemoteAddr()

	// Increase request count
	count := &(h.requests)
	atomic.AddUint64(count, 1)

	if ctx.IsGet() {

		url := ctx.URI()
		filename := string(ctx.Path())

		// GET parameters
		query := url.QueryArgs()

		qw, _ := query.GetUint("width")
		qh, _ := query.GetUint("height")
		qm := string(query.Peek("mode"))
		uqw := uint(qw)
		uqh := uint(qh)
		puqw := &uqw
		puqh := &uqh
		puqm := &qm
		h.RetrieveImage(ctx, filename, puqw, puqh, puqm)

	} else if ctx.IsPost() {
		// POST is currently unused so we can use this for testing
		h.RetrieveHello(ctx)
	}
}

// Respond to POST message by saying Hello
func (h MyHandler) RetrieveHello(ctx *fasthttp.RequestCtx) {
	_, err := ctx.WriteString("Hello world!")
	if err != nil {
		log.Println(err)
	}
}

// Write image by filename into ResponseWriter with the
// desired width and height being pointed to. If there
// are no desired width or height, that parameter is nil.
func (h *MyHandler) RetrieveImage(ctx *fasthttp.RequestCtx,
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
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		ctx.WriteString("Image not found!")
		logging.Debug(err)
		return
	}
	ctx.Write(blob)
}

// Create a new fasthttp server and configure it.
// This does not run the server however.
func NewServer(maxMemory uint64) (*fasthttp.Server, *MyHandler, net.Listener) {

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
	server := &fasthttp.Server{
		Handler: 				handler.ServeHTTP,
		ReadTimeout:		8 * time.Second,
		WriteTimeout:		8 * time.Second,
		ReadBufferSize:	1 << 20,			// this also limits max header size
	}
	ln, _ := net.Listen("tcp", ":8005")
	return server, handler, ln
}

func main() {

	// CLI arguments
	mem := flag.Uint64("m", 500*1024*1024, "Sets the maximum memory to be used for caching images in bytes. Does not account for memory consumption of other things.")
	flag.Parse()

	server, handler, ln := NewServer(*mem)

	log.Println("Server starting...")
	handler.started = time.Now()
	//err := server.ListenAndServe(":8005")
	err := server.Serve(ln)
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
