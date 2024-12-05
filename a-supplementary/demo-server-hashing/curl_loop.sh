#!/bin/bash

# URL to send requests to
url="http://localhost:8080/hash"

# Send 5 requests to the URL
for i in {1..100000}
do
   curl $url
done