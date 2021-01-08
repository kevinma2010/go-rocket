package core

import "errors"

var ErrAppPathEmpty = errors.New("project(<app_path, a>) path can not be empty")
var ErrNotInGoPath = errors.New("the project directory must be in the $GOPATH/src")
