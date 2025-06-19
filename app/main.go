package main

import (
	"github.com/codecrafters-io/shell-starter-go/app/logger"
	"github.com/codecrafters-io/shell-starter-go/app/shell"
)

func main() {
	logger.InitLogger()
	shell := shell.New()
	shell.Run()
}
