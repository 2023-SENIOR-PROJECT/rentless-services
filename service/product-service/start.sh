#!/bin/bash

/root/main &
/root/client &
# go run main.go &
# go run client.go &
wait -n
exit &?
