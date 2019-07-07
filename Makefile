init:
	# install fdocker
	go install .
	# create network(docker0)
	fdocker network create --driver bridge --subnet 192.168.0.0/16 test

run:
	fdocker run --ti --name test1 --net test -m 128m --images $(shell dirname $(shell pwd)/$(lastword $(MAKEFILE_LIST)))/images/ busybox /bin/sh

clean:
	# stop fdocker
	fdocker stop test
	# remove fdocker
	fdocker rm test
	# clean network
	fdocker network remove test
