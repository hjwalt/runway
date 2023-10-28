#!/bin/bash

openssl req -x509 \
            -sha256 \
            -days 356 \
            -nodes \
            -newkey rsa:2048 \
            -subj "/CN=hjwalt.me/C=SG/L=Singapore" \
            -keyout root.key -out root.crt 


openssl genrsa -out server.key 2048

openssl req -new -key server.key -out server.csr -config csr.conf

openssl x509 -req \
    -in server.csr \
    -CA root.crt -CAkey root.key \
    -CAcreateserial -out server.crt \
    -days 365 \
    -sha256 -extfile cert.conf