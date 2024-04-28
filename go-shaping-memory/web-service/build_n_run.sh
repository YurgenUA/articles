#!/bin/bash

go build;
if [ $? -ne 0 ]; \
then \
    echo "build failed :("; \
    exit $?; \
fi

echo "Running...."
./http-server