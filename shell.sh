#!/bin/bash

docker run -ti\
	-v $PWD:/go/src/github.com/tadasv/gohlc\
	--workdir='/go/src/github.com/tadasv/gohlc'\
	golang:1.5.1 bash
