build :
	go build -o ./bin/slblog ./cmd/slblog

run :
	go run ./cmd/slblog

install :
	cp ./bin/slblog /usr/local/bin
