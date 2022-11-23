#!/bin/sh

sqlite3 shapeless-blog.db "ALTER TABLE posts ADD preview TEXT NOT NULL DEFAULT '';" ".exit"
