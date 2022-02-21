#!/bin/bash

declare -a proto_paths=("../../auth/rpc/auth" "../../things/rpc/auth" "../../users/rpc/auth")
start_dir="$(pwd)"

for i in "${proto_paths[@]}"
do
  protoc --go_out=$i --go_opt=paths=source_relative \
  --go-grpc_out=$i --go-grpc_opt=paths=source_relative \
  auth.proto 
done