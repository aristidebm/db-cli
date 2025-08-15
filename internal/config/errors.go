package config

import (
	"errors"
)

var InvalidSource = errors.New("invalid source")
var InvalidClient = errors.New("invalid client")
var UnsupportedDriver = errors.New("unsupported driver")
