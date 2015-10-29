package main

import (
	"gopkg.in/tylerb/graceful.v1"
	"log"
	"net/http"
	"server/cache"
	"strconv"
	"sync/atomic"
	"time"
)

// MyHandler type is used to encompass HandlerFunc interface.
// In the future this type will probably contain pointers to
// services provided by this program (image cache).
type MyHandler struct {

	// Additional output (debug)
	verbose bool

	// Service started
	started time.Time

	// Request count (statistics)
	requests uint64

	// ImageCache
	images *cache.ImageCache
}

// ServeHTTP is called whenever there is a new request.
// This is quite similar to JavaEE Servlet interface.
// TODO: Check that ServeHTTP is called inside a goroutine?
func (self *MyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method := (*request).Method

	// In the future we can use requester can detect request spammers!
	// requester := (*request).RemoteAddr

	// Increase request count
	count := &((*self).requests)
	atomic.AddUint64(count, 1)

	// SPLIT on method
	if method == "GET" {
		// GET an image by name
		req := (*request).URL

		// Get the filename
		filename := (*req).Path

		// Get parameters
		query := req.Query()

		// Extract width and height if needed
		var width *int = nil
		var height *int = nil
		if len(query) > 0 {
			strw, ok := query["width"]
			if ok && len(strw) > 0 {
				intw, err := strconv.Atoi(strw[0])
				if err == nil {
					width = new(int)
					*width = intw
				}
				// For now, silent error if !ok
			}
			strh, ok := query["height"]
			if ok && len(strh) > 0 {
				inth, err := strconv.Atoi(strh[0])
				if err == nil {
					height = new(int)
					*height = inth
				}
				// For now, silent error if !ok
			}
		}

		// Get the image
		(*self).RetrieveImage(&writer, filename, width, height)

	} else if method == "POST" {
		// POST is currently unused so we can use this for testing
		(*self).RetrieveHello(&writer)
	}
}

// Respond to POST message by saying Hello
func (*MyHandler) RetrieveHello(writer *http.ResponseWriter) {
	result := "Hello world!"
	_, err := (*writer).Write([]byte(result))
	if err != nil {
		log.Println(err)
	}
}

// Write image by filename into ResponseWriter with the
// desired width and height being pointed to. If there
// are no desired width or height, that parameter is nil.
func (self *MyHandler) RetrieveImage(writer *http.ResponseWriter,
	filename string,
	width *int,
	height *int) {

	// TODO: filename must not be interpret as "absolute"
	// implement a type that will abstract away the filesystem.
	if (*self).verbose {
		log.Println("Find: " + filename)
	}

	// Get cache
	bank := (*self).images

	// Load the image
	blob, err := (*bank).GetImage(filename, width, height)
	if err != nil {
		// TODO:
		// Classify different possible errors more but make sure
		// no "internal" information is leaked.
		(*writer).WriteHeader(http.StatusNotFound)
		(*writer).Write([]byte("Image not found!"))
		return
	}
	(*writer).Write(blob)
}

// Create a new graceful server and configure it.
// This does not run the server however.
func NewServer() (*graceful.Server, *MyHandler) {
	handler := &MyHandler{
		// TODO: Set this to true while debugging (preprosessor possible?)
		verbose: true,

		// Initialize
		requests: 0,

		// No cache (later sprints)
		images: &Cacheless{},
	}
	server := &graceful.Server{
		Timeout: 8 * time.Second,
		Server: &http.Server{
			Addr:           ":8005",
			Handler:        handler,
			ReadTimeout:    8 * time.Second,
			WriteTimeout:   8 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
	return server, handler
}

// The entry point for the program.
// This is obviously not exported.
func main() {

	// Allocate server
	server, handler := NewServer()

	// Run the server
	begin := time.Now()
	(*handler).started = begin
	err := server.ListenAndServe()
	end := time.Now()

	// Get number of requests
	requests := strconv.FormatUint((*handler).requests, 10)

	// Calculate the elapsed time
	duration := end.Sub(begin)
	log.Println("Server requests: " + requests)
	log.Println("Server uptime: " + duration.String())

	// Log errors
	if err != nil {
		log.Fatal(err)
	}
}
