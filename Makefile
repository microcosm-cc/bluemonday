# Targets:
#
#   all:          Builds the code locally after testing
#
#   fmt:          Formats the source files
#   build:        Builds the code locally
#   vet:          Vets the code
#   lint:         Runs lint over the code (you do not need to fix everything)
#   test:         Runs the tests
#   cover:        Gives you the URL to a nice test coverage report
#   clean:        Deletes the built file (if it exists)
#
#   install:      Builds, tests and installs the code locally

# The first target is always the default action if `make` is called without args
# We clean, build and install into $GOPATH so that it can just be run
all: clean fmt vet test install

fmt:
	gofmt -w ./$*

build: clean
	GOOS=linux GOARCH=amd64 go build

vet:
	go tool vet *.go

lint:
	golint *.go

test:
	go test -v ./...

cover:
	go test -coverprofile=coverage.out && \
	go tool cover -html=coverage.out && rm coverage.out

clean:
	find $(GOPATH)/pkg/*/github.com/microcosm-cc -name bluemonday.a -delete

install: clean
	go install
