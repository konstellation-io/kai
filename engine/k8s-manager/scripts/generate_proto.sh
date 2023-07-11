#!/bin/sh

PROTO_PATH=./internal/infrastructure/grpc/proto

for PROTO_FILE in internal/infrastructure/grpc/proto/**/*.proto; do
  protoc -I=$PROTO_PATH/versionpb/ \
  --go_out=$PROTO_PATH/versionpb --go_opt=paths=source_relative \
  --go-grpc_out=$PROTO_PATH/versionpb --go-grpc_opt=paths=source_relative \
  $PROTO_FILE
done

if [ -d "../admin-api" ]; then
  for PB_FILE in $(find $PROTO_PATH -iname "*.pb.go"); do
    DST=$(basename "$(dirname "$PB_FILE")")
    mkdir -p ../admin-api/adapter/service/proto/${DST}
    cp "$PB_FILE" "../admin-api/adapter/service/proto/${DST}/"
  done
fi
