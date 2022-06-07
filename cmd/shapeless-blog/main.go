package main

import (
	"flag"
	"os"
)

func main() {
	var port int
	var shapelessBlogDir string
	flag.StringVar(&shapelessBlogDir, "dir", os.Getenv("SHAPELESS_BLOG_DIR"), "shapeless-blog directory")
	flag.IntVar(&port, "port", 9398, "shapeless-blog port")
	flag.Parse()

}
