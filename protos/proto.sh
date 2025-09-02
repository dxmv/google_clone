#!/bin/bash

# generate protos for search service (Go)
export PATH="$PATH:$(go env GOPATH)/bin"
protoc --proto_path=. --go_out=../services/search/pb --go_opt=paths=source_relative --go-grpc_out=../services/search/pb --go-grpc_opt=paths=source_relative search.proto

# generate protos for query-api service (Python)
cd ../services/query-api 
source .venv/bin/activate && python -m grpc_tools.protoc -I../../protos --python_out=pb --pyi_out=pb --grpc_python_out=pb ../../protos/search.proto

# fix Python imports for relative imports
sed -i '' 's/^import search_pb2 as search__pb2$/from . import search_pb2 as search__pb2/' pb/search_pb2_grpc.py

echo "Protobuf files generated successfully!"
echo "Go files: ../services/search/pb/"
echo "Python files: pb/"