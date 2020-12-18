#!/bin/bash

ADDR=http://127.0.0.1:8081/v1

curl -i -X POST -F 'user_id=1' -F 'currency=USD' -F 'amount=1000' ${ADDR}/wallets/

curl -i -X POST -F 'user_id=2' -F 'currency=USD' -F 'amount=1000' ${ADDR}/wallets/

curl -i -X POST -F 'user_id=3' -F 'currency=USD' -F 'amount=1000' ${ADDR}/wallets/