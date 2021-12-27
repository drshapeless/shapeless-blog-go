build :
	go build -o ./bin/slblog ./main

run :
	go run ./main/

install :
	cp slblog /usr/local/bin
