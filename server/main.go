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

	config      config.Conf
	operator    operator.Operator
	imageSource ops.ImageSource
	watermarker ops.Watermarker
}

var defaults = config.ConfValues{
	MinHeight:  200,
	MinWidth:   200,
	MaxHeight:  5000,
	MaxWidth:   5000,
	AddMark:    false,
	Imgpath:    "",
	Tokens:     1,
	Vertical:   0.0,
	Horizontal: 1.0,
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
		operations, format, err, invalid := ParseURI(url, h.imageSource, h.watermarker, h.config)
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
			ctx.Write(blob)
			logging.Debug("Blob returned")
		}

	} else if ctx.IsPost() {
		// POST is currently unused so we can use this for testing
		h.RetrieveHello(ctx)
		logging.Debug("Post request received")
	}
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
	logging.Debug("Reading server config")
	//setting default values

	minHeight, err := conf.GetInt("watermark", "minheight")
	if err != nil {
		log.Printf("Error reading config size minimum height restriction, defaulting to %v\n", defaults.MinHeight)
		minHeight = defaults.MinHeight
	}

	minWidth, err := conf.GetInt("watermark", "minwidth")
	if err != nil {
		log.Printf("Error reading config size minimum width restriction, defaulting to %v\n", defaults.MinWidth)
		minWidth = defaults.MinWidth
	}

	maxHeight, err := conf.GetInt("watermark", "maxheight")
	if err != nil {
		log.Printf("Error reading config size maximum height restriction, defaulting to %v\n", defaults.MaxHeight)
		maxHeight = defaults.MaxHeight
	}

	maxWidth, err := conf.GetInt("watermark", "maxwidth")
	if err != nil {
		log.Printf("Error reading config size maximum width restriction, defaulting to %v\n", defaults.MaxWidth)
		maxWidth = defaults.MaxWidth
	}

	addMark, err := conf.GetBool("watermark", "addmark")
	if err != nil {
		log.Println("Error reading config addmark value, defaulting to false")
	}

	imgpath, err := conf.GetString("watermark", "path")
	if err != nil && addMark == true {
		log.Println("Error reading path for watermark image, disabling watermarking")
		addMark = false
	}

	ver, err := conf.GetFloat64("watermark", "vertical")
	if err != nil {
		log.Printf("Error reading config vertical alignment, defaulting to %v\n", defaults.Vertical)
		ver = defaults.Vertical
	}

	hor, err := conf.GetFloat64("watermark", "horizontal")
	if err != nil {
		log.Printf("Error reading config horizontal alignment, defaulting to %v\n", defaults.Horizontal)
		hor = defaults.Horizontal
	}

	tokens, err := conf.GetInt("server", "concurrency")
	if err != nil {
		log.Printf("Error reading config concurrency value, defaulting to %v\n", defaults.Tokens)
		tokens = defaults.Tokens
	}

	watermarker, err := ops.MakeWatermarker(imgpath, hor, ver, maxWidth, minWidth, maxHeight, minHeight, addMark)
	if err != nil {
		log.Printf("Error creating watermarker: %v\n", err.Error())
	}

	// Configure handler
	logging.Debug("Configuring handler")
	handler := &MyHandler{
		requests:    0,
		config:      conf,
		imageSource: imageSource,
		operator:    operator.MakeDefault(maxMemory, "/tmp/RICdiskcache", tokens),
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

	cpath := flag.String("c", "config.ini", "Sets the configuration .ini file used.")
	flag.Parse()
	// CLI arguments

	conf, err := config.ReadConfig(*cpath)
	if err != nil {
		log.Printf("Error while reading config at %v: %v, using default values\n", *cpath, err.Error())
	}
	def, err := conf.GetUint64("server", "memory")
	if err != nil {
		log.Printf("Error reading config memory value, defaulting to %v\n", defaults.Mem)
		def = defaults.Mem
	}
	mem := flag.Uint64("m", def, "Sets the maximum memory to be used for caching images in bytes. Does not account for memory consumption of other things.")
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
