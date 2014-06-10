# Targets:
#
#   all:          Builds the code locally after testing
#
#   fmt:          Formats the source files
#   build:        Builds the code locally
#   vet:          Vets the code
#   lint:         Runs lint over the code (you do not need to fix everything)
#   test:         Runs the tests
#   clean:        Deletes the built file (if it exists)
#
#   install:      Builds, tests and installs the code locally

# Sub-directories containing code to be vetted or linted
CODE = *.go

# The first target is always the default action if `make` is called without args
# We clean, build and install into $GOPATH so that it can just be run
all: clean fmt vet test install

fmt:
	gofmt -w ./$*

build: clean
	GOOS=linux GOARCH=amd64 go build

vet:
	go tool vet $(CODE)

lint:
	golint $(CODE)

test:
	go test -v -cover ./...

clean:
	find $(GOPATH)/bin -name bluemonday -delete
	find . -name bluemonday -delete

install: clean
	go install
