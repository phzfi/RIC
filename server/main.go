package main

import (
	"flag"
	"fmt"
	"gopkg.in/gographics/imagick.v2/imagick"
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/logging"
	"github.com/phzfi/RIC/server/images"
	"github.com/phzfi/RIC/server/ops"
	"github.com/phzfi/RIC/server/config"
	"github.com/valyala/fasthttp"
	"log"
	"net"
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

	config			config.Conf
	operator    cache.Operator
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

		url := ctx.URI()
		operations, err := ParseURI(url, h.imageSource, h.watermarker, h.config)
		if err != nil {
			ctx.NotFound()
			logging.Debug(err)
			return
		}
		blob, err := h.operator.GetBlob(operations...)
		if err != nil {
			ctx.NotFound()
			logging.Debug(err)
		} else {
			// ping blob to avoid reading whole image into memory
			// and set right content-type response header
			img := images.NewImage()
			img.PingImageBlob(blob)
			ctx.SetContentType("image/" + img.GetImageFormat())
			img.Destroy()

			ctx.Write(blob)
			logging.Debug("Blob returned")
		}

	} else if ctx.IsPost() {
		// POST is currently unused so we can use this for testing
		h.RetrieveHello(ctx)
		logging.Debug("Post request received")
	}
}

func getParams(a *fasthttp.Args) (w *uint, h *uint, m string) {
	qw, e := a.GetUint("width")
	if e == nil {
		uqw := uint(qw)
		w = &uqw
	}
	qh, e := a.GetUint("height")
	if e == nil {
		uqh := uint(qh)
		h = &uqh
	}

	m = string(a.Peek("mode"))
	return
}

// Respond to POST message by saying Hello
func (h MyHandler) RetrieveHello(ctx *fasthttp.RequestCtx) {
	_, err := ctx.WriteString("Hello world!")
	if err != nil {
		log.Println(err)
	}
}

// Create a new fasthttp server and configure it.
// This does not run the server however.
func NewServer(port int, maxMemory uint64, conf config.Conf) (*fasthttp.Server, *MyHandler, net.Listener) {
	logging.Debug("Creating server")
	imageSource := ops.MakeImageSource()
	// Add roots
	// TODO: This must be externalized outside the source code.
	logging.Debug("Adding roots")
	if imageSource.AddRoot("/var/www") != nil {
		log.Fatal("Root not added /var/www")
	}

	if imageSource.AddRoot(".") != nil {
		log.Println("Root not added .")
	}
	imgpath, err := conf.GetString("watermark", "path")
	watermarker := ops.MakeWatermarker(imgpath)

	if err != nil {
		log.Fatal("Error creating listener:" + err.Error())
	}

	// Configure handler
	logging.Debug("Configuring handler")
	handler := &MyHandler{
		requests:    0,
		config: conf,
		imageSource: imageSource,
		watermarker: watermarker,
		operator:    cache.MakeOperator(maxMemory),
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

	cpath := flag.String("c", "config.ini", "Sets the configuration .ini file used.")
	flag.Parse()
	// CLI arguments

	conf, err := config.ReadConfig(*cpath)
	if err != nil {
		log.Fatal("Error while reading config at " + *cpath + ": " + err.Error())
	}

	def, err := conf.GetUint64("server", "memory")
	if err != nil {
		def = 512*1024*1024
	}
	mem := flag.Uint64("m", def, "Sets the maximum memory to be used for caching images in bytes. Does not account for memory consumption of other things.")
	flag.Parse()

	imagick.Initialize()
	defer imagick.Terminate()

	log.Println("Server starting...")
	logging.Debug("Debug enabled")

	server, handler, ln := NewServer(8005, *mem, conf)
	handler.started = time.Now()
	err = server.Serve(ln)
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
