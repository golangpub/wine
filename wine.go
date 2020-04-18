package wine

import (
	"github.com/golangpub/log"
	"github.com/golangpub/wine/internal/path"
)

var logger *log.Logger

func init() {
	logger = log.Default().Derive("Wine")
	logger.SetFlags(log.LstdFlags - log.Lfunction - log.Lshortfile)
	path.SetLogger(logger)
}

func Logger() *log.Logger {
	return logger
}
