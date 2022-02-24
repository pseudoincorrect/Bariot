#!/bin/bash

declare -a proto_paths=("auth")

start_dir="$(pwd)"

for i in "${proto_paths[@]}"
do
  cd "$start_dir/$i"

  protoc --go_out=. --go_opt=paths=source_relative \
  --go-grpc_out=. --go-grpc_opt=paths=source_relative \
  $i.proto 
done