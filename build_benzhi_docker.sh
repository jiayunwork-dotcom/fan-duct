#!/bin/bash
set -e

IMAGE_NAME=${1:-my-project}
DOCKER_PLATFORM=${2:-linux/amd64}

docker build --platform $DOCKER_PLATFORM -f benzhi.Dockerfile -t $IMAGE_NAME .

echo ""
echo "Docker image '$IMAGE_NAME' built successfully."
echo ""
echo "Next steps:"
echo "  docker run -d -P --name fan-duct-b14 $IMAGE_NAME:latest"
echo "  curl -s http://127.0.0.1:\$(docker port fan-duct-b14 8080 | cut -d: -f2)/api/example"
echo "  docker rm -f fan-duct-b14"
echo ""
