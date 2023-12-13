# RIC
Responsive image cache


### Development environment

## Vagrant
```bash
# update vagrant stuff
./vagrant_up.sh

# start virtual machine
vagrant up

# start RIC server
vagrant ssh

cd /vagrant/server

go build

./do_run.sh

# check that script was run, to detach `ctrl + a` and `d`
screen -list

# shutdown vagrant
./vagrant_down.sh
```

## Docker
```bash
# build and start docker container
./up.sh

# shutdown docker container
./down.sh
```


Test that server returns test images:
`http://localhost:8005/01.jpg`

### Accepted RIC query parameters

* width
* height
* mode: fit, liquid, crop
* format: All that imagemagic supports
* watermark: text
* url: webroot url
