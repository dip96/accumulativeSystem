# shellcheck disable=SC2164
cd ./cmd/gophermart
rm -rf gophermart
go build -o gophermart *.go
