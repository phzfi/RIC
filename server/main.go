package main

import (
	"flag"
	"fmt"
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/images"
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

		if string(ctx.Path()) == "/health" || string(ctx.Path()) == "/healthz" {
			ctx.SetContentType("application/json")
			ctx.WriteString(`{"status":"ok"}`)
			return
		}

		url := ctx.URI()
		operations, format, _, err, invalid := ParseURI(url, h.imageSource, h.watermarker)
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
func (h *MyHandler) RetrieveHello(ctx *fasthttp.RequestCtx) {
	_, err := ctx.WriteString("Hello world!")
	if err != nil {
		log.Println(err)
	}
}

func buildHybridCache(cfg *config.ConfValues, memoryLimit uint64) cache.Cacher {
	caches := []cache.Cacher{
		cache.NewCache(cache.NewLRU(), memoryLimit),
	}

	if cfg.Cache.DiskPath != "" {
		diskMaxBytes := cfg.Cache.DiskMaxMB * 1024 * 1024
		disk := cache.NewDiskCache(cfg.Cache.DiskPath, diskMaxBytes, cache.NewLRU())
		if disk != nil {
			caches = append(caches, disk)
			logging.Debugf("Added disk cache: %s (max %d MB)", cfg.Cache.DiskPath, cfg.Cache.DiskMaxMB)
		}
	}

	if cfg.Cache.S3Enabled {
		s3Cache, err := cache.NewS3Cache(cfg.Cache, cache.NewLRU())
		if err == nil {
			caches = append(caches, s3Cache)
			logging.Debugf("Added S3 cache: bucket=%s, prefix=%s", cfg.Cache.S3Bucket, cfg.Cache.S3Prefix)
		} else {
			log.Printf("Warning: S3 cache unavailable: %v", err)
		}
	}

	logging.Debugf("Hybrid cache created with %d tiers", len(caches))
	return cache.HybridCache(caches)
}

// Create a new fasthttp server and configure it.
// This does not run the server however.
func NewServer(port int, maxMemory uint64, conf *config.ConfValues) (*fasthttp.Server, *MyHandler, net.Listener) {
	logging.Debug("Creating server")

	var imageSource ops.ImageSource

	if conf.ImageSource.S3Enabled {
		s3client, err := images.NewS3Client(conf.ImageSource)
		if err != nil {
			log.Printf("Warning: S3 image source unavailable: %v", err)
			imageSource = ops.MakeImageSource()
		} else {
			imageSource = ops.MakeImageSourceWithS3(s3client)
			logging.Debugf("S3 image source configured: bucket=%s, prefix=%s", conf.ImageSource.S3Bucket, conf.ImageSource.S3Prefix)
		}
	} else {
		imageSource = ops.MakeImageSource()
	}

	// Add local file roots
	logging.Debug("Adding roots")
	if imageSource.AddRoot("/var/www") != nil {
		log.Fatal("Root not added /var/www")
	}

	if imageSource.AddRoot(".") != nil {
		log.Println("Root not added .")
	}

	if imageSource.AddRoot("/testimages/server") != nil {
		log.Println("Root not added /testimages/server")
	}
	logging.Debug("Reading server config")
	//setting default values

	watermarker, err := ops.MakeWatermarker(conf.Watermark)
	if err != nil {
		log.Printf("Error creating watermarker: %v\n", err.Error())
	}

	hc := buildHybridCache(conf, maxMemory)

	// Configure handler
	logging.Debug("Configuring handler")
	handler := &MyHandler{
		requests:    0,
		imageSource: imageSource,
		operator:    operator.Make(hc, conf.Server.Tokens),
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

	conf := config.ReadConfig()

	mem := flag.Uint64("m", conf.Server.Memory, "Sets the maximum memory to be used for caching images in bytes. Does not account for memory consumption of other things.")
	flag.Parse()
	imagick.Initialize()
	defer imagick.Terminate()

	log.Println("Server starting...")
	logging.Debug("Debug enabled")

	server, handler, ln := NewServer(8005, *mem, conf)
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
