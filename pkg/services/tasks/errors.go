package tasks

import "errors"

var errRegisterAfterServe = errors.New("tasks: RegisterService called after Serve")
