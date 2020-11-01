package main

import "github.com/task-executor/pkg/controllers/build"

func main() {
	bc := build.NewBuildController()
	bc.Start()
}
