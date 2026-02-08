#!/bin/bash
# generate-proto.sh - Script to generate Go code from Protocol Buffers

echo "Generating Protocol Buffers Go code..."

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Please download from: https://github.com/protocolbuffers/protobuf/releases"
    echo "For Windows: download protoc-<version>-win64.zip and add to PATH"
    exit 1
fi

# Check if Go plugins are installed
if ! command -v protoc-gen-go &> /dev/null; then
    echo "Installing protoc-gen-go..."
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo "Installing protoc-gen-go-grpc..."
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Get the directory of this script
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# Create directory if not exists
mkdir -p shared/proto

# Generate Go code
echo "Generating Go code from model.proto..."
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    shared/proto/model.proto

if [ $? -eq 0 ]; then
    echo "✅ Successfully generated protobuf code!"
    echo "Generated files:"
    echo "  - shared/proto/model.pb.go"
    echo "  - shared/proto/model_grpc.pb.go"
else
    echo "❌ Failed to generate protobuf code"
    exit 1
fi
