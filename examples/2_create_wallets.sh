#!/bin/bash

ADDR=http://127.0.0.1:8081/v1

curl -i -X POST -F 'user_id=32' -F 'currency=USD' -F 'value=1000' ${ADDR}/wallets/
echo "\n"

curl -i -X POST -F 'user_id=33' -F 'currency=USD' -F 'value=1000' ${ADDR}/wallets/
echo "\n"