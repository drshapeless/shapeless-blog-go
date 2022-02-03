build :
	go build -o ./bin/slblog ./cmd/slblog

run :
	go run ./main/

install :
	cp ./bin/slblog /usr/local/bin
