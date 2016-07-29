# RIC
Responsive image cache


### Development environment
```bash
# start virtual machine
vagrant up

# start RIC server
vagrant ssh

cd /vagrant/server

go build

./do_run.sh

# check that script was run, to detach `ctrl + a` and `d`
screen -list
```

Test that server returns test images:
`http://localhost:8005/01.jpg`

### Accepted RIC query parameters

* width
* height
* mode: fit, liquid, crop
* format: All that imagemagic supports
* watermark: text