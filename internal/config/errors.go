package config

import (
	"errors"
)

var InvalidSource = errors.New("invalid source")
var InvalidClient = errors.New("invalid client")
var UnsupportedScheme = errors.New("unsupported scheme")
