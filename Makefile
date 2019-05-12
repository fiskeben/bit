.PHONY: test clean

bit:main.go
	go build .

test:
	go test .

clean:
	rm ./bit
