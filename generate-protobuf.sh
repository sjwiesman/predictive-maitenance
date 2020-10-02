#!/bin/bash

if ! command -v protoc &> /dev/null
then
    echo "protoc could not be found"
    exit
fi

protoc \
  --go_out=digital-twins \
  --go_out=simulator \
  --python_out=model \
  protobuf/cloudfun.proto

mv model/protobuf/cloudfun_pb2.py model/cloudfun_pb2.py