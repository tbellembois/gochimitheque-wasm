build:
	GOARCH=wasm GOOS=js go build -o wasm
	gzip -9 -v -c wasm > ../gochimitheque/wasm/wasm.gz
	rm wasm
	