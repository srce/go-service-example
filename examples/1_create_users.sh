#!/bin/bash

ADDR=http://127.0.0.1:8081/v1

curl -i -X POST -F 'name=JEFF BEZOS' -F 'email=jeff@amazon.com' ${ADDR}/users/
echo "\n"

curl -i -X POST -F 'name=BILL GATES' -F 'email=bill@microsoft.com' ${ADDR}/users/
echo "\n"