package main

/*
 * Enter directory examples/websites/hello
 * go run ./main.go
 */

import (
	"github.com/golangpub/wine"
)

func main() {
	s := wine.NewServer()
	s.StaticDir("/", "./html")
	s.Run(":8000")
}
