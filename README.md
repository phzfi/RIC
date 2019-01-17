# RIC
Responsive image cache

TODO: Update development environment setup



### Development environment

Create folder structure `ric_project_gopath/src/github.com/phzfi`
```bash
mkdir -p ric_project_gopath/src/github.com/phzfi
cd ric_project_gopath/src/github.com/phzfi

git clone <REPO_URL>

```

Run docker-composer and login
```bash
docker-compose up --build --force-recreate
docker exec -i -t ric_dev /bin/bash
```

Inside docker container
```bash
scripts/provision/docker_install_go_dependencies.sh

go build

./server

```

#Configuration

Default location of ric config is `/etc/ric/ric_config.ini.`
In configuration file, remote server whitelist configuration file path must be defined.
Default location is `/mnt/config/host_whitelist.ini`


Test that server returns status page:
`http://localhost:8005/status`

### Accepted RIC query parameters

* width
* height
* mode: fit, liquid, crop
* format: All that imagemagic supports
* watermark: text
* url: webroot url



#Compile and run delve
Ide must be configured to respond to connection
```bash

go build -tags debug -v -gcflags "all=-N -l" && /root/go/bin/dlv --headless --listen=:40000 --api-version=2 exec ./server

```
