#!make
include .envrc

current_time = $(shell date --iso-8601=seconds)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = "-s -X 'main.buildTime=${current_time}' -X 'main.version=${git_description}'"

build:
	go build -ldflags=${linker_flags} -o=./bin/shapeless-blog ./cmd/shapeless-blog

run:
	go run ./cmd/shapeless-blog

test:
	./cmd/shapeless-blog -path=${SHAPELESS_BLOG_DB_PATH} -secret=${SHAPELESS_BLOG_SECRET}

clean:
	rm ./bin/shapeless-blog

install:
	cp ./bin/shapeless-blog /usr/local/bin/shapeless-blog

uninstall:
	rm /usr/local/bin/shapeless-blog
