package main

import (
	"flag"
	"fmt"
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/operator"
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
	"gopkg.in/gographics/imagick.v2/imagick"
	"log"
	"net"
	"os"
	"path/filepath"
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

	operator    operator.Operator
	imageSource ops.ImageSource
	watermarker ops.Watermarker
}

// ServeHTTP is called whenever there is a new request.
// This is quite similar to JavaEE Servlet interface.
func (h *MyHandler) ServeHTTP(ctx *fasthttp.RequestCtx) {

	// In the future we can use requester can detect request spammers!
	// requester := ctx.RemoteAddr()

	// Increase request count
	count := &(h.requests)
	atomic.AddUint64(count, 1)

	if ctx.IsGet() {

		// Special case for status check.
		// TODO: Consider implementing routing?
		path := string(ctx.Path())
		if path == "/status" {
			_, err := ctx.WriteString("OK")
			if err != nil{
				ctx.Error("Failed to write output", 500)
			}
			return
		}

		uri := ctx.URI()
		operations, format, err, invalid := ParseURI(uri, h.imageSource, h.watermarker)
		if err != nil {
			ctx.NotFound()
			logging.Debug(err)
			return
		}
		if invalid != nil {
			ctx.Error(invalid.Error(), 400)
			return
		}
		blob, err := h.operator.GetBlob(operations...)
		if err != nil {
			ctx.NotFound()
			logging.Debug(err)
		} else {
			ctx.SetContentType("image/" + format)

			length, err := ctx.Write(blob)
			if err != nil {
				ctx.Error("Failed to write output", 500)
				return
			}
			logging.Debug(fmt.Sprintf("Blob returned with length: %d", length))
		}

	} else if ctx.IsPost() {
		// POST is currently unused so we can use this for testing
		h.RetrieveHello(ctx)
		logging.Debug("Post request received")
	}
}

// Respond to POST message by saying Hello
func (h MyHandler) RetrieveHello(ctx *fasthttp.RequestCtx) {
	_, err := ctx.WriteString("OK")
	if err != nil {
		log.Println(err)
	}
}

// Create a new fasthttp server and configure it.
// This does not run the server however.
func NewServer(port int, maxMemory uint64, conf *config.ConfValues) (*fasthttp.Server, *MyHandler, net.Listener) {
	logging.Debug("Creating server")
	imageSource := ops.MakeImageSource()

	// Add roots
	logging.Debug("Adding roots")
	if conf.Server.ImageFolder != "" {
		if imageSource.AddRoot(conf.Server.ImageFolder) != nil {
			log.Fatal(fmt.Sprintf("Root not added %s", conf.Server.ImageFolder))
		}
	}

	if imageSource.AddRoot(".") != nil {
		log.Println("Root not added .")
	}
	logging.Debug("Reading server config")
	//setting default values

	watermarker, err := ops.MakeWatermarker(conf.Watermark)
	if err != nil {
		log.Printf("Error creating watermarker: %v\n", err.Error())
	}

	// Configure handler
	logging.Debug("Configuring handler")
	handler := &MyHandler{
		requests:    0,
		imageSource: imageSource,
		operator:    operator.MakeDefault(maxMemory, "/tmp/RICdiskcache", conf.Server.Tokens),
		watermarker: watermarker,
	}

	// Configure server
	server := &fasthttp.Server{
		Handler: handler.ServeHTTP,
	}

	logging.Debug("Beginning to listen")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatal("Error creating listener:" + err.Error())
	}
	logging.Debug("Server ready")
	return server, handler, ln
}

func main() {

	cpath := locateConfig()

	logging.Debug(fmt.Sprintf("Loading config from %s", *cpath))
	flag.Parse()
	// CLI arguments

	conf := config.ReadConfig(*cpath)

	mem := flag.Uint64("m", conf.Server.Memory, "Sets the maximum memory to be used for caching images in bytes. Does not account for memory consumption of other things.")
	imagick.Initialize()
	defer imagick.Terminate()

	log.Println(fmt.Sprintf("Server starting. Listening to port %d...", conf.Server.Port ))
	logging.Debug("Debug enabled")

	server, handler, ln := NewServer( conf.Server.Port, *mem, conf)

	handler.started = time.Now()
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

func locateConfig() *string {
	// First check directory from where binary was launched
	currentDirectory, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	if _, err := os.Stat(currentDirectory + "/ric_config.ini"); os.IsNotExist(err) {
		// No try default installation location
		currentDirectory = getBinaryFileDirectory()
	}

	if _, err := os.Stat(currentDirectory + "/ric_config.ini"); os.IsNotExist(err) {
		log.Fatal("Could not load configuration file")
	}
	cpath := flag.String("c", currentDirectory + "/ric_config.ini", "Sets the configuration .ini file used.")

	return cpath
}

func getBinaryFileDirectory() string {

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exPath := filepath.Dir(ex)
	return exPath
}
