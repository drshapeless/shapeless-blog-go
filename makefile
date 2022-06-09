#!make
include .envrc

build:
	go build -o ./bin/shapeless-blog ./cmd/shapeless-blog

run:
	go run ./cmd/shapeless-blog

test:
	./cmd/shapeless-blog -path=${SHAPELESS_BLOG_DB_PATH} -secret=${SHAPELESS_BLOG_SECRET}

clean:
	rm ./bin/shapeless-blog
