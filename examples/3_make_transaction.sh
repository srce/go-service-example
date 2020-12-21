#!/bin/bash

ADDR=http://127.0.0.1:8081/v1

curl -i -X POST -F 'sender_id=2' -F 'beneficiary_id=3' -F 'amount=100.00' -F 'currency=USD' ${ADDR}/transactions/