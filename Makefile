all:
	echo 'Provide a target: publish clean'

fmt:
	find src/ -name '*.go' -exec go fmt {} ';'

build: fmt
	gb build all

publish: build
	PORT=9896 ./bin/publish

test:
	gb test -v

clean:
	rm -rf bin/ pkg/

.PHONY: publish
