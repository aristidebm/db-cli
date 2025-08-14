package client

import (
	"errors"
)

var URLParseError = errors.New("cannot not parse connection url")
var UnsupportedCommand = errors.New("unsupported command")
