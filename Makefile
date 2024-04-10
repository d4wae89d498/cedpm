all:
	cd src/cli && go build -o cedpm

test:
	export PATH=$$PATH:$$(pwd)/src/cli && cd examples && cedpm test
