install:
	go clean
	go install

run:
	go run lhc.go

clean:
	go clean
	rm -f lhc
