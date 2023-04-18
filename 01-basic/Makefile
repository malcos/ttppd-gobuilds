# use always go mod init before working on any go projects, go modules are not perfect but better than the old strategy
export GO111MODULE=on

# temporal file for storing coverage output, needed later for visualization
COVERAGE_FILE=coverage.out

# profit from go modules and get the project name directly from the module file
PRJNAME := $(lastword $(subst /, ,$(word 2, $(file < go.mod))))

$(PRJNAME):
	go build -o $@ ./cmd

clean:
	rm -f $(PRJNAME) $(COVERAGE_FILE)

test:
	go test -count=1 ./...

test.coverage:
	go test -cover -coverprofile $(COVERAGE_FILE) -count=1 ./...

test.coverage.view: test.coverage
	go tool cover -html=$(COVERAGE_FILE)
