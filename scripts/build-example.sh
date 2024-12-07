TAG="nickname/image"
VERSION="latest"

docker rmi ${TAG}:${VERSION}
docker build -t ${TAG}:${VERSION} .