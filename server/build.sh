set -e
echo "building go server for Linux..."
GOOS=linux go build
docker build -t bond00729/game-api .
docker push bond00729/game-api
go clean