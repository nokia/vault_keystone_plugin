EXTERNAL_TOOLS=\
               github.com/mitchellh/gox \
               github.com/kardianos/govendor

GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# bootstrap the build by downloading additional tools
bootstrap:
	@for tool in  $(EXTERNAL_TOOLS) ; do \
		echo "Installing/Updating $$tool" ; \
		go get -u $$tool;
	done

fmt:
	gofmt -w $(GOFMT_FILES)

dev:
