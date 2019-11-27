GOOS=linux
GOARCH=amd64
BUILDPATH=$(CURDIR)
BINPATH=$(BUILDPATH)/bin
EXENAME=apibin

clean:
	@sudo rm -r /tmp/pgdata || true
	@sudo rm -r bin/ || true
	@sudo rm -r log/ || true

build: clean
	@if [ ! -d $(BINPATH) ] ; then mkdir -p $(BINPATH) ; fi
	@CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $(BINPATH)/$(EXENAME) -a -installsuffix . || (echo "build failed $$?"; exit 1)
	@echo 'Build suceeded... done'

