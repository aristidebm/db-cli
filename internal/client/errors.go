package client

import (
	"errors"
)

var UnsupportedCommand = errors.New("unsupported command")
var UnsupportedFormat = errors.New("unsupported format")
