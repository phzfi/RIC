package main

import (
	"github.com/phzfi/RIC/server/cache"
	"github.com/phzfi/RIC/server/ops"
)


func SetupOperatorSource() (operator cache.Operator, src ops.ImageSource) {
	operator = cache.MakeOperator(512 * 1024 * 1024)
	src = ops.MakeImageSource()
	src.AddRoot("./")
	return
}


