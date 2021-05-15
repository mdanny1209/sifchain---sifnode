#!/bin/bash

set -e

if  [[ ! -f  ~/.docker/config.json || $(cat ~/.docker/config.json  | jq '.auths["ghcr.io"].auth') == 'null' ]]; then
  echo "In order to run this script and push a new container to the github registry you need to create a personal access token and use it to login to ghcr with docker"
  echo ""
  echo "  echo \$MY_PAT | docker login ghcr.io -u USERNAME --password-stdin"
  echo ""
  echo "For more information see https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry#authenticating-to-the-container-registry"
  echo ""
  echo "Create a personal access token and log into docker using the above link then try running this script again.s"
  exit 1
fi

if [[ ! -z "$(git status --porcelain --untracked-files=no)" ]]; then 
  echo "Git workspace must be clean to save git commit hash"
  exit 1
fi

echo "Github Registry Login found."
echo "Building new container..."

TAG=$(git rev-parse HEAD)
IMAGE_NAME=ghcr.io/sifchain/sifnode/ui-stack:$TAG

# Assume this script is run from `./ui`
ROOT=$(pwd)/..

echo "New image name: $IMAGE_NAME"

# Using buildkit to take advantage of local dockerignore files
export DOCKER_BUILDKIT=1

cd $ROOT && docker build -f ./ui/scripts/stack.Dockerfile -t $IMAGE_NAME .

docker push $IMAGE_NAME

echo $IMAGE_NAME > $ROOT/ui/scripts/latest

echo "Commit the ./ui/scripts/latest file to git"