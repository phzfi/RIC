# Responsive Image Cache

## 1. Project Description

### 1.1. Business Vision

Responsive Image Cache (RIC) is an open source (image) CDN server optimized for high speed caching and delivery of images, in exact optimized pixel sizes of the clients.

Instead of resizing small, medium, large images of the original (e.g. RAW media bank files), RIC can deliver exact sizes to each viewport and device. This
- makes images load faster
- reduces wasted bandwidth and too large image sizes
- reduces rendering time on client by skipping the need to resize the image by the browser

Nice features include
- GPU acceleration
- Liquid rescale
- text watermarks

For client side see src/riclib.js for example usage by Javascript.

See also RIC Wordpress plugin https://github.com/phzfi/ric-wordpress

### 1.2. Task Management

### 1.3. Personas

### 1.4. Use Cases

### 1.5. Non-Functional Requirements

## 2. Architecture

### 2.1. Technologies

Mostly build in Golang.

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

Add here examples and hints of good ways how to code the project. Convert the silent knowledge as tacit knowledge here.

## 3. Development Environment
Note! PHZ Coding Convention: name this environment as dev.
Note! However, please use the default files for dev env, such as docker-compose.yml (instead of docker-compose.dev.yml).

### 3.1. Prerequisites

### 3.2. Start the Application

Start docker env

    ./up.sh

Tear down

    ./down.sh

Status

    ./status.sh

### 3.3. Access the Application

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

Test that server returns test images:

    http://localhost:8105/01.jpg

#### Accepted RIC HTTP query parameters

* width: int in px
* height: int in px
* mode: fit, liquid, crop
* format: All that Imagemagic supports
* watermark: text
* url: webroot url of source images

For example http://localhost:8105/01.jpg?width=200&height=300&mode=liquid&watermark=PHZ.fi

### 3.4. Run Tests

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

