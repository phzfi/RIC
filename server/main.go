package main

import (
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/logging"
	"gopkg.in/tylerb/graceful.v1"
	"io"
	"log"
	"net/http"
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

	images cache.ImageCache
}

// ServeHTTP is called whenever there is a new request.
// This is quite similar to JavaEE Servlet interface.
// TODO: Check that ServeHTTP is called inside a goroutine?
func (h *MyHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method := request.Method

	// In the future we can use requester can detect request spammers!
	// requester := request.RemoteAddr

	// Increase request count
	count := &(h.requests)
	atomic.AddUint64(count, 1)

	// SPLIT on method
	if method == "GET" {
		// GET an image by name
		req := request.URL

		filename := req.Path

		// Get parameters
		query := req.Query()

		// Extract width and height if needed
		var width, height *uint

		if len(query) > 0 {
			strw, ok := query["width"]
			if ok && len(strw) > 0 {
				intw, err := strconv.ParseUint(strw[0], 10, 32)
				if err == nil {
					*width = uint(intw)
				}
				// For now, silent error if !ok
			}
			strh, ok := query["height"]
			if ok && len(strh) > 0 {
				inth, err := strconv.ParseUint(strh[0], 10, 32)
				if err == nil {
					*height = uint(inth)
				}
				// For now, silent error if !ok
			}
		}

		h.RetrieveImage(writer, filename, width, height)

	} else if method == "POST" {
		// POST is currently unused so we can use this for testing
		h.RetrieveHello(writer)
	}
}

// Respond to POST message by saying Hello
func (h MyHandler) RetrieveHello(writer http.ResponseWriter) {
	result := "Hello world!"
	_, err := io.WriteString(writer, result)
	if err != nil {
		log.Println(err)
	}
}

// Write image by filename into ResponseWriter with the
// desired width and height being pointed to. If there
// are no desired width or height, that parameter is nil.
func (h *MyHandler) RetrieveImage(writer http.ResponseWriter,
	filename string,
	width *uint,
	height *uint) {

	// TODO: filename must not be interpret as "absolute"
	// implement a type that will abstract away the filesystem.
	logging.Debug("Find: " + filename)

	// Get cache
	bank := h.images

	// Load the image
	blob, err := bank.GetImage(filename, width, height)
	if err != nil {
		// TODO:
		// Classify different possible errors more but make sure
		// no "internal" information is leaked.
		writer.WriteHeader(http.StatusNotFound)
		io.WriteString(writer, "Image not found!")
		logging.Debug(err)
		return
	}
	writer.Write(blob)
}

// Create a new graceful server and configure it.
// This does not run the server however.
func NewServer() (*graceful.Server, *MyHandler) {
	// No cache (will be implemented in later sprints)
	cacher := &cache.Cacheless{}

	// Add roots
	// TODO: This must be externalized outside the source code.
	if cacher.AddRoot("/var/www") != nil {
		log.Fatal("Root not added /var/www")
	}

	if cacher.AddRoot(".") != nil {
		log.Println("Root not added .")
	}

	// Root count
	if len(cacher.Roots) != 2 {
		log.Fatal("All roots not added")
	}

	// Configure handler
	handler := &MyHandler{
		requests: 0,
		images:   cacher,
	}

	// Configure server
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
