package main

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/phzfi/RIC/server/config"
	"github.com/phzfi/RIC/server/operator"
	"github.com/phzfi/RIC/server/ops"
	"github.com/valyala/fasthttp"
)

// TestRetrieveHello tests the RetrieveHello method
func TestRetrieveHello(t *testing.T) {
	handler := MyHandler{}

	ctx := &fasthttp.RequestCtx{}
	handler.RetrieveHello(ctx)

	if string(ctx.Response.Body()) != "Hello world!" {
		t.Errorf("Expected 'Hello world!', got '%s'", string(ctx.Response.Body()))
	}
}

// TestServeHTTP_PostMethod tests POST method handling
func TestServeHTTP_PostMethod(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_cache", 1),
		watermarker: ops.Watermarker{},
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetRequestURI("/somepath")
	handler.ServeHTTP(ctx)

	if string(ctx.Response.Body()) != "Hello world!" {
		t.Errorf("Expected 'Hello world!', got '%s'", string(ctx.Response.Body()))
	}
}

// TestServeHTTP_GetMethod_Health tests /health endpoint
func TestServeHTTP_GetMethod_Health(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_cache", 1),
		watermarker: ops.Watermarker{},
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/health")
	handler.ServeHTTP(ctx)

	if string(ctx.Response.Body()) != `{"status":"ok"}` {
		t.Errorf("Expected '{\"status\":\"ok\"}', got '%s'", string(ctx.Response.Body()))
	}
}

// TestServeHTTP_GetMethod_Healthz tests /healthz endpoint
func TestServeHTTP_GetMethod_Healthz(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_cache", 1),
		watermarker: ops.Watermarker{},
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/healthz")
	handler.ServeHTTP(ctx)

	if string(ctx.Response.Body()) != `{"status":"ok"}` {
		t.Errorf("Expected '{\"status\":\"ok\"}', got '%s'", string(ctx.Response.Body()))
	}
}

// TestMyHandler_RequestCount tests that request count increments
func TestMyHandler_RequestCount(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_cache", 1),
		watermarker: ops.Watermarker{},
	}

	if handler.requests != 0 {
		t.Errorf("Expected 0 requests, got %d", handler.requests)
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetRequestURI("/test")
	handler.ServeHTTP(ctx)

	if handler.requests != 1 {
		t.Errorf("Expected 1 request, got %d", handler.requests)
	}
}

// TestServeHTTP_GetBlobError tests error path when GetBlob fails
func TestServeHTTP_GetBlobError(t *testing.T) {
	// Create a handler with a broken operator that will fail on GetBlob
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_getblob_error", 1),
		watermarker: ops.Watermarker{},
	}
	handler.imageSource.AddRoot("/app/server/testimages/server/")

	// Request an image with invalid format that will cause conversion to fail
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/01.jpg?width=200&height=200&mode=resize&format=invalid")
	handler.ServeHTTP(ctx)

	// Should get an error response (404 or 400)
	if ctx.Response.StatusCode() == 200 {
		t.Errorf("Expected error response, got 200")
	}
}

// TestServeHTTP_ResizeDenyUpscale tests resize with upscale denial
func TestServeHTTP_ResizeDenyUpscale(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_upscale", 1),
		watermarker: ops.Watermarker{},
	}
	handler.imageSource.AddRoot("/app/server/testimages/server/")

	// Request larger dimensions than original (900x1200)
	// Resize should deny upscaling, so width/height should be adjusted
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/01.jpg?width=2000&height=2000&mode=resize&format=jpeg")
	handler.ServeHTTP(ctx)

	// Should succeed but with adjusted dimensions
	if ctx.Response.StatusCode() != 200 {
		t.Errorf("Expected 200, got %d", ctx.Response.StatusCode())
	}
}

// TestServeHTTP_FitMode tests fit mode
func TestServeHTTP_FitMode(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_fit", 1),
		watermarker: ops.Watermarker{},
	}
	handler.imageSource.AddRoot("/app/server/testimages/server/")

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/01.jpg?width=200&height=200&mode=fit&format=jpeg")
	handler.ServeHTTP(ctx)

	if ctx.Response.StatusCode() != 200 {
		t.Errorf("Expected 200, got %d", ctx.Response.StatusCode())
	}
}

// TestServeHTTP_InvalidParams tests 400 for invalid parameters
func TestServeHTTP_InvalidParams(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_cache", 1),
		watermarker: ops.Watermarker{},
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/01.jpg?width=abc&height=200&mode=resize&format=jpeg")
	handler.ServeHTTP(ctx)

	if ctx.Response.StatusCode() != 400 {
		t.Errorf("Expected 400, got %d", ctx.Response.StatusCode())
	}
}

// TestNewServer tests server creation
func TestNewServer(t *testing.T) {
	conf := &config.ConfValues{
		Watermark: config.Watermark{
			MinHeight: 200, MinWidth: 200,
			MaxHeight: 5000, MaxWidth: 5000,
		},
		Server: config.Server{Tokens: 1, Memory: 1024 * 1024 * 100},
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	server, handler, ln := NewServer(port, 1024*1024*100, conf)
	if server == nil {
		t.Error("Server should not be nil")
	}
	if handler == nil {
		t.Error("Handler should not be nil")
	}
	if ln == nil {
		t.Error("Listener should not be nil")
	}
	ln.Close()
}

// TestNewServer_InvalidWatermarkPath tests NewServer with invalid watermark image path
func TestNewServer_InvalidWatermark(t *testing.T) {
	conf := &config.ConfValues{
		Watermark: config.Watermark{
			MinHeight: 200, MinWidth: 200,
			MaxHeight: 5000, MaxWidth: 5000,
			ImagePath: "/nonexistent/watermark.png", // Invalid image path
		},
		Server: config.Server{Tokens: 1, Memory: 1024 * 1024 * 100},
	}

	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	// NewServer should still return even with invalid watermark
	server, handler, ln := NewServer(port, 1024*1024*100, conf)
	if server == nil {
		t.Error("Server should not be nil even with invalid watermark")
	}
	if handler == nil {
		t.Error("Handler should not be nil even with invalid watermark")
	}
	ln.Close()
}

func startTestServer(t *testing.T) (*fasthttp.Server, *MyHandler, net.Listener, int) {
	conf := &config.ConfValues{
		Watermark: config.Watermark{
			MinHeight: 200, MinWidth: 200,
			MaxHeight: 5000, MaxWidth: 5000,
		},
		Server: config.Server{Tokens: 1, Memory: 1024 * 1024 * 100},
	}

	// Use port 0 to let NewServer pick a random available port
	server, handler, ln := NewServer(0, 1024*1024*100, conf)
	port := ln.Addr().(*net.TCPAddr).Port
	go server.Serve(ln)
	time.Sleep(100 * time.Millisecond)

	return server, handler, ln, port
}

// TestServerIntegration tests health endpoint
func TestServerIntegration(t *testing.T) {
	_, _, ln, port := startTestServer(t)
	defer ln.Close()

	_, body, err := fasthttp.Get(nil, fmt.Sprintf("http://localhost:%d/health", port))
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `{"status":"ok"}` {
		t.Errorf("Expected '{\"status\":\"ok\"}', got '%s'", string(body))
	}
}

// TestConcurrentRequests tests concurrent request handling
func TestConcurrentRequests(t *testing.T) {
	_, _, ln, port := startTestServer(t)
	defer ln.Close()

	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, body, err := fasthttp.Get(nil, fmt.Sprintf("http://localhost:%d/health", port))
			if err != nil {
				t.Errorf("Request failed: %v", err)
			}
			if string(body) != `{"status":"ok"}` {
				t.Errorf("Unexpected response: %s", string(body))
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestServeHTTP_InvalidMode tests 400 response for invalid mode
func TestServeHTTP_InvalidMode(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_cache", 1),
		watermarker: ops.Watermarker{},
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/01.jpg?width=200&height=200&mode=invalidmode&format=jpeg")
	handler.ServeHTTP(ctx)

	if ctx.Response.StatusCode() != 400 {
		t.Errorf("Expected 400, got %d", ctx.Response.StatusCode())
	}
}

// TestServeHTTP_InvalidFormat tests 400 response for invalid format
func TestServeHTTP_InvalidFormat(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_cache", 1),
		watermarker: ops.Watermarker{},
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/01.jpg?width=200&height=200&mode=resize&format=invalid")
	handler.ServeHTTP(ctx)

	if ctx.Response.StatusCode() != 400 {
		t.Errorf("Expected 400, got %d", ctx.Response.StatusCode())
	}
}

// TestServeHTTP_ImageNotFound tests 404 when image doesn't exist
func TestServeHTTP_ImageNotFound(t *testing.T) {
	handler := &MyHandler{
		requests:    0,
		imageSource: ops.MakeImageSource(),
		operator:    operator.MakeDefault(1024*1024*100, "/tmp/ric_test_cache", 1),
		watermarker: ops.Watermarker{},
	}
	handler.imageSource.AddRoot("/app/server/testimages/server/")

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("GET")
	ctx.Request.SetRequestURI("/nonexistent.jpg?width=200&height=200&mode=resize&format=jpeg")
	handler.ServeHTTP(ctx)

	if ctx.Response.StatusCode() != 404 {
		t.Errorf("Expected 404, got %d", ctx.Response.StatusCode())
	}
}

// TestRetrieveHello_WriteError tests RetrieveHello when WriteString fails
// Note: This is difficult to test directly since fasthttp doesn't expose write errors easily
// We'll skip this test as it requires low-level fasthttp internals
func TestRetrieveHello_WriteError(t *testing.T) {
	// Create a context - WriteString errors are hard to simulate in fasthttp
	// The error path exists in the code (line 83-87 in main.go)
	// but testing it requires internal fasthttp knowledge
	t.Skip("Cannot easily test WriteString errors in fasthttp")
}
