#! /bin/bash

openssl genrsa -out docs/private.pem 3072
openssl rsa -pubout -in docs/private.pem -out docs/public.pem