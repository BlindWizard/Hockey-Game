#!/bin/sh

go version
go build -o /home/go/server -buildvcs=false 
cd /home/go
./server