all:
	cd src/cli && go build -o ../../bin/cedpm

test:
	cd examples && cedpm d d
