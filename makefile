# Makefile example
build:
	go build -o bin/gocr-lib .

run:
	./bin/gocr-lib $(filter-out $@, $(MAKECMDGOALS))

clean:
	rm -rf ./bin/

test:
	go test -v ./...