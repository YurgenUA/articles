#!/bin/bash

go build
if [ $? -ne 0 ]; 
then 
    echo "build failed :("; 
    exit 1; \
fi

echo "Running...."
./http-server