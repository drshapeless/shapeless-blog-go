#!/bin/sh

sudo rc-service shapeless-blog stop

sudo cp ~/shapeless-blog /usr/local/bin/shapeless-blog

sudo rc-service shapeless-blog start
