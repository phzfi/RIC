# Responsive Image Cache

## 1. Project Description

### 1.1. Business Vision

Responsive Image Cache (RIC) is an open source (image) CDN server optimized for high speed caching and delivery of images, in exact optimized pixel sizes of the clients.

Instead of resizing small, medium, large images of the original (e.g. RAW media bank files), RIC can deliver exact sizes to each viewport and device. This
- makes images load faster
- reduces wasted bandwidth and too large image sizes
- reduces rendering time on client by skipping the need to resize the image by the browser

**Key Features:**
- GPU acceleration via ImageMagick
- Liquid rescale (content-aware scaling / seam carving)
- Text and image watermarks
- Multiple resize modes (resize, fit, crop, cropmid, liquid)
- Format conversion (jpeg, gif, webp, bmp, png, tiff)
- Hybrid caching (LRU memory + disk cache)
- Concurrent request handling with configurable tokens

For client side usage see `src/riclib.js` for example usage by JavaScript.

See also RIC WordPress plugin https://github.com/phzfi/ric-wordpress

### 1.2. Task Management

Source code can be found from Github https://github.com/phzfi/RIC . Please feel free to contribute!

Licensed under permissive open source MIT -license. See LICENSE.

### 1.3. Personas

### 1.4. Use Cases

**Basic Image Resizing:**
```
http://localhost:8105/01.jpg?width=200
http://localhost:8105/01.jpg?width=200&height=300
http://localhost:8105/01.jpg?width=200&format=webp
```

**Resize Modes:**
| Mode | Description | Example |
|------|-------------|---------|
| `resize` | Scale to fit (default) | `?width=200&height=300&mode=resize` |
| `fit` | Fit within bounds, maintain aspect | `?width=200&height=300&mode=fit` |
| `crop` | Crop from top-left corner | `?width=200&height=300&mode=crop&cropx=50&cropy=50` |
| `cropmid` | Crop from center | `?width=200&height=300&mode=cropmid` |
| `liquid` | Content-aware scaling | `?width=200&height=300&mode=liquid` |

**Watermarks:**
```
# Enable watermark from URL parameter
http://localhost:8105/01.jpg?width=200&watermark=true

# Disable watermark
http://localhost:8105/01.jpg?width=200&watermark=false

# Add text watermark
http://localhost:8105/01.jpg?width=200&watermark=MyText
```

**External Image Sources:**
```
http://localhost:8105/image.jpg?width=200&url=https://example.com/images/
```

#### Image Sources

RIC supports multiple types of image sources:

**Local filesystem roots:**
- Configured in `server/main.go` (AddRoot calls)
- Default roots: `/var/www`, `.`, `/testimages/server`
- Images are loaded from `filepath.Join(root, filename)`

**Web roots (URL-based):**
- Roots prefixed with `http://` or `https://`
- Example: `https://example.com/images/`
- Images are fetched via HTTP request

**URL parameter:**
```
?url=https://example.com/images/
```
- Dynamically adds a web root for the request
- Must be prefixed with `http://` or `https://`

### 1.5. Non-Functional Requirements

## 2. Architecture

### 2.1. Technologies

Mostly built in Golang with the following components:

- **fasthttp** - High-performance HTTP server
- **ImageMagick (imagick)** - Image processing with GPU acceleration
- **Custom hybrid cache** - LRU memory + disk cache

All PHZ Full Stack -projects should encapsulate all environments by virtualization:

Dev
* (Vagrant/Virtualbox) - deprecated
* Docker Compose/Docker
** Images available at https://hub.docker.com/repository/docker/phzfi/ric

CI
* use dev -env on ci.in.phz.fi + Jenkins executors running Docker Compose
* Jenkins
* do not pin the projects down on any individual executor, but set up the builds so that they can be run on any executor machine

Staging
* PHZ Docker Swarm
* Kubernetes

Production
* PHZ Docker Swarm (internal projects only)
* Kubernetes
* AWS
* or any other environment of your wish

### 2.2. Naming, Terms and Key Concepts

Environments and the configs should be named as
* dev: docker-compose.yml (i.e. use the default names for dev env), but .env.dev
* (ci): use the dev -env on CI
* stg: docker-compose.stg.yml, .env.stg
* prod: docker-compose.prod.yml, .env.prod

### 2.3. Coding Convention

Directory structure
* docs/ for documentation
* etc/ for nginx, ssh etc configs. Can be cp -pr etc/ /etc to the virtual machine during provisioning and matches the os directory structure
* results/ test results
* reports/ for e.g. code coverage reports
* src/ for source code
** Note! Source code should be placed under a single folder (src) that can be mounted over Docker -volume or Vagrant -shared folder inside the virtual machine so that node_modules or vendor directory are not on the shared folder. See https://wiki.phz.fi/Docker and https://wiki.phz.fi/Vagrant for further details how to circumvent the problems.
* tests/ for tests

### 2.4. Development Guide

**Component Architecture:**
```
Client -> fasthttp Server -> ParseURI -> Operator -> HybridCache
                              |
                              +-> ImageSource (loads from roots: /var/www, .)
                              |
                              +-> Image Operations (Resize, Crop, LiquidRescale, Watermark, Convert)
                              |
                              +-> ImageMagick (GPU accelerated)
```

**Key Components:**
- `server/main.go` - Server entry point, handles HTTP requests via MyHandler
- `server/urlparser.go` - Parses query parameters, creates operation chain
- `server/operator/operator.go` - Orchestrates image processing with caching and concurrency
- `server/cache/hybridcache.go` - LRU memory cache + disk cache (4GB disk, 1GB chunks)
- `server/ops/roots.go` - Manages image file roots and loading
- `server/ops/watermarker.go` - Image and text watermark operations

#### Cache Types

RIC uses a hybrid caching system with configurable eviction policies:

**Cache Policies:**

| Policy | Description | File |
|--------|-------------|------|
| `LRU` (default) | Least Recently Used - evicts least accessed items first | `cache/lru.go` |
| `FIFO` | First In First Out - evicts oldest items first | `cache/fifo.go` |

**Cache Layers:**

| Layer | Description | Default Config |
|-------|-------------|----------------|
| Memory Cache | LRU-backed in-memory store | 2GB max (`SERVER_MEMORY`) |
| Disk Cache | Persistent on-disk storage | 4GB max, folder `/tmp/RICdiskcache` |

**How HybridCache Works:**
```
Request -> Memory Cache (LRU) -> Disk Cache (LRU) -> Image Processing
                 |                      |
                 +-- Cache Hit! --------+
                 |
                 +-- Cache Miss -> Process -> Store in both layers
```

**Eviction Behavior:**
- Memory cache: LRU policy evicts least recently visited items when memory limit reached
- Disk cache: Stores base64-encoded blobs in files, uses LRU policy
- When image found in disk cache but not memory, it's promoted to memory cache

**Configuration:**
- `SERVER_MEMORY` - Memory cache limit in bytes (default: 2147483648 = 2GB)
- Disk cache folder is set in `server/main.go` (default: `/tmp/RICdiskcache`)

## 3. Development Environment
Note! PHZ Coding Convention: name this environment as dev.
Note! However, please use the default files for dev env, such as docker-compose.yml (instead of docker-compose.dev.yml).

### 3.1. Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- ImageMagick with Wand library
- NVIDIA GPU and nvidia-container-toolkit (for GPU acceleration)

### 3.2. Start the Application

Start docker env

    ./up.sh

Tear down

    ./down.sh

Status

    ./status.sh

### 3.3. Access the Application

Test that server returns test images:

    http://localhost:8105/01.jpg

Health check endpoints:

    http://localhost:8105/health
    http://localhost:8105/healthz

#### Configuration

All configuration is done via environment variables in `.env.dev` (development), `.env.stg` (staging), and `.env.prod` (production). The legacy `config.ini` files have been removed.

**Watermark settings:**

| Variable | Description | Default |
|---|---|---|
| `WATERMARK_PATH` | Path to watermark image | `""` |
| `WATERMARK_HORIZONTAL` | Horizontal position (0.0-1.0) | `1.0` |
| `WATERMARK_VERTICAL` | Vertical position (0.0-1.0) | `0.0` |
| `WATERMARK_MAXWIDTH` | Max image width to apply watermark | `5000` |
| `WATERMARK_MINWIDTH` | Min image width to apply watermark | `200` |
| `WATERMARK_MAXHEIGHT` | Max image height to apply watermark | `5000` |
| `WATERMARK_MINHEIGHT` | Min image height to apply watermark | `200` |
| `WATERMARK_ADDMARK` | Enable watermark by default | `false` |
| `WATERMARK_TEXT` | Default watermark text (used when AddMark=true and no URL param) | `""` |
| `WATERMARK_FORCEMARK` | Force watermark (overrides URL params) | `""` |

**Server settings:**

| Variable | Description | Default |
|---|---|---|
| `SERVER_TOKENS` | Concurrency tokens | `1` |
| `SERVER_MEMORY` | Memory limit in bytes | `2147483648` (2GB) |

#### Accepted RIC HTTP Query Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `width` | int | Target width in pixels |
| `height` | int | Target height in pixels |
| `mode` | string | Resize mode: resize, fit, crop, cropmid, liquid |
| `format` | string | Output format: jpeg, png, webp, gif, bmp, tiff |
| `watermark` | string | Watermark control: true, false, or custom text |
| `url` | string | External image source URL root |
| `cropx` | int | X offset for crop mode |
| `cropy` | int | Y offset for crop mode |

**Example requests:**
```
http://localhost:8105/01.jpg?width=200
http://localhost:8105/01.jpg?width=200&height=300&mode=liquid
http://localhost:8105/01.jpg?width=200&watermark=PHZ.fi
```

#### JavaScript Client (riclib.js)

```html
<script>
window.RICConfig = {
    server_path: 'http://localhost:8105',
    maxres: 1920,
    quality: 85
};
</script>
<script src="src/riclib.js"></script>
```

The library automatically processes all `<img>` tags on page load, converting them to use the RIC server for optimized delivery.

### 3.4. Run Tests

```bash
./test.sh
```

View coverage report:
```bash
./coverage.sh
# Open reports/coverage/coverage.html
```

### 3.5. IDE Setup and Debugging

### 3.6. Version Control

### 3.7. Databases and Migrations

### 3.8. Continuous Integration

## 4. Staging Environment
Note! PHZ Coding Convention: name this environment as stg.

### 4.1. Access

### 4.2. Deployment

### 4.3. Smoke Tests

#### 4.3.1. Automated Test Cases

#### 4.3.2. Manual Test Cases

### 4.4. Rollback

### 4.5. Logs

### 4.6. Monitoring

## 5. Production Environment
Note! PHZ Coding Convention: name this environment as prod.

### 5.1. Access

### 5.2. Deployment

### 5.3. Smoke Tests

#### 5.3.1. Automated Test Cases

#### 5.3.2. Manual Test Cases

### 5.4. Rollback

### 5.5. Logs

### 5.6. Monitoring

## 6. Operating Manual

### 6.1. Scheduled Jobs

### 6.2. Manual Processes

### 6.3. Devops Life Cycle Management Plan

Add here known information of estimates how fast the chosen technologies and versions will be deprecated.

## 7. Problems

### 7.1. Environments

### 7.2. Coding

### 7.3. Dependencies

Add here TODO and blockers that you have found related to upgrading to newer versions.
List the library/framework/service, version, and then the error message.