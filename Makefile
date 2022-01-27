test:
	go test -gcflags=-l -p 1 -timeout 0 -count 1 -covermode=atomic -coverprofile=coverage.txt -v ./...
	go tool cover -func=coverage.txt
lint:
	golines . -w
	golangci-lint run
	
build:
	go build
install:
	go get -v -t -d ./...
race_test:
	go test -gcflags=-l -p 1 -timeout 0 -count 1 -race -v ./...

protolint_install:
	GO111MODULE=on GOPATH=~/go go get -u -v github.com/yoheimuta/protolint/cmd/protolint@v0.27.0
protolint:
	PATH=${PATH}:${GOPATH}/bin/:~/go/bin protolint lint -config_path=./protolint.yaml -fix ./aea/mail ./packages/fetchai/protocols
protolint_install_win:
	powershell -command '$$env:GO111MODULE="on"; go get -u -v github.com/yoheimuta/protolint/cmd/protolint@v0.27.0'
protolint_win:
	protolint lint -config_path=./protolint.yaml -fix ./aea/mail ./packages/fetchai/protocols
