package io

import "github.com/golangpub/log"

var logger = log.Default()

func SetLogger(l *log.Logger) {
	logger = l
}
