package todo

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskAlredyExists = errors.New("task alredy exists")
