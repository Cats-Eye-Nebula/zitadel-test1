#! /bin/sh

set -eux

echo "Generate grpc"

#TODO: find a way to generate swagger and authoption to the correct package without mv
mkdir $GOPATH/src/github.com/caos/zitadel/swagger
mkdir /proto/output
move() {
  mv /proto/output/zitadel/$1*.swagger.json $GOPATH/src/github.com/caos/zitadel/swagger/
  mv /proto/output/zitadel/$1*.go $GOPATH/src/github.com/caos/zitadel/pkg/grpc/$1/
}

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  /proto/include/zitadel/message.proto

protoc \
  -I=/proto/include \
  --go_out ${GOPATH}/src \
  --go-grpc_out ${GOPATH}/src \
  --grpc-gateway_out $GOPATH/src/github.com/caos/zitadel/pkg/grpc \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out $GOPATH/src/github.com/caos/zitadel/swagger \
  --openapiv2_opt logtostderr=true \ 
  --authoption_out=/proto/output \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/include/zitadel/admin.proto

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --grpc-gateway_out $GOPATH/src/github.com/caos/zitadel/pkg/grpc \
  --grpc-gateway_opt logtostderr=true \
  --grpc-gateway_opt allow_delete_body=true \
  --openapiv2_out $GOPATH/src/github.com/caos/zitadel/swagger \
  --openapiv2_opt logtostderr=true \ 
  --authoption_out=/proto/output \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/include/zitadel/management.proto

protoc \
  -I=/proto/include \
  --go_out $GOPATH/src \
  --go-grpc_out $GOPATH/src \
  --grpc-gateway_out $GOPATH/src/github.com/caos/zitadel/pkg/grpc \
  --grpc-gateway_opt logtostderr=true \
  --openapiv2_out $GOPATH/src/github.com/caos/zitadel/swagger \
  --openapiv2_opt logtostderr=true \ 
  --authoption_out=/proto/output \
  --validate_out=lang=go:${GOPATH}/src \
  /proto/include/zitadel/auth.proto

move "admin"
move "auth"
move "management"

echo "done generating grpc"