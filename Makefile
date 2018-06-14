
.PHONY: setup
setup-linter:
	go get -u gopkg.in/alecthomas/gometalinter.v2
	go get -u golang.org/x/tools/cmd/gotype
	go get -u github.com/fzipp/gocyclo
	go get -u golang.org/x/lint/golint
	go get -u github.com/opennota/check/cmd/aligncheck
	go get -u github.com/opennota/check/cmd/structcheck
	go get -u github.com/opennota/check/cmd/varcheck
	go get -u github.com/kisielk/errcheck
	go get -u honnef.co/go/tools/cmd/megacheck
	go get -u github.com/mibk/dupl
	go get -u github.com/gordonklaus/ineffassign
	go get -u mvdan.cc/interfacer
	go get -u github.com/mdempsky/unconvert
	go get -u github.com/jgautheron/goconst/cmd/goconst
	go get -u github.com/GoASTScanner/gas/cmd/gas/...

.PHONY: setup
setup: setup-linter
	go get -u golang.org/x/tools/cmd/goimports
	go get -t -v ./...

.PHONY: fmt
fmt:
	goimports -w=true -d .

.PHONY: lint
lint:
	gometalinter.v2 ./...

.PHONY: test
test:
	go test -cover -race ./...
