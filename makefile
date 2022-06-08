#!make
include .envrc

build:
	go build -o ./bin/shapeless-blog ./cmd/shapeless-blog

run:
	go run ./cmd/shapeless-blog

clean:
	rm ./bin/shapeless-blog
