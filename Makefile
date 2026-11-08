.PHONY: test example memory-test benchmark

test:
	go test ./...

example:
	go run ./cmd/kmeans -n 100000 -k 1000

memory-test:
	go test -tags mlpack_memtest -run TestNative -v .

benchmark:
	go build -o ./kmeans-mesure ./cmd/kmeans
	python3 tools/benchmark.py --binary ./kmeans-mesure --output docs/validation/mesures.json
