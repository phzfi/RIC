# RIC
Responsive image cache

Project structure attempts to follow the structure defined at: `https://github.com/golang-standards/project-layout`



### Development environment


```bash
git clone https://github.com/phzfi/RIC.git

```

Run docker-composer and login
```bash
docker-compose up
docker exec -i -t ric_dev /bin/bash
```

Inside docker container
```bash
/ric/scripts/build_dev.sh

/ric/server/server

```

#Configuration

Default location of ric config is `/etc/ric/ric_config.ini`
In configuration file, remote server whitelist configuration file path must be defined.
Default location is `/var/lib/ric/config/host_whitelist.ini`


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
cd /root/go/src/github.com/phzfi/RIC/server
go build -tags debug -v -gcflags "all=-N -l" && /root/go/bin/dlv --headless --listen=:40000 --api-version=2 exec ./server

```
