.PHONY : test
test :
	go run ./cmd/shapeless-blog

build :
	go build -o ./bin/shapeless-blog ./cmd/shapeless-blog
